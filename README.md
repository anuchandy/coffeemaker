# coffeemaker
UncleBob's Mark IV CoffeeMaker [Implementation in Go]

### The Mark IV Special Coffee Maker [Problem statement - taken from Uncle Bob's [article] (http://objectmentor.com/resources/articles/CoffeeMaker.pdf)]


The Mark IV Special makes up to 12 cups of coffee at a time. The user places a filter in the filter holder, fills the filter with coffee grounds, and slides the filter holder into its receptacle. The user then pours up to 12 cups of water into the water strainer and presses the Brew button. The water is heated until boiling. The pressure of the evolving steam forces the water to be sprayed over the coffee grounds, and coffee drips through the filter into the pot. The pot is kept warm for extended periods by a warmer plate, which turns on only if coffee is in the pot. If the pot is removed from the warmer plate while water is being sprayed over the grounds, the flow of water is stopped so that brewed coffee does not spill on the warmer plate. The following hardware needs to be monitored or controlled:

* The heating element for the boiler. It can be turned on or off.
* The heating element for the warmer plate. It can be turned on or off.
* The sensor for the warmer plate. It has three states: warmerEmpty, potEmpty, potNotEmpty.
* A sensor for the boiler, which determines whether water is present. It has two states: boilerEmpty or boilerNotEmpty.
* The Brew button. This momentary button starts the brewing cycle. It has an indicator that lights up when the brewing cycle is over and the coffee is ready.
* A pressure-relief valve that opens to reduce the pressure in the boiler. The drop in pressure stops the flow of water to the filter. The value can be opened or closed.

The hardware for the Mark IV has been designed and is currently under development. The hardware engineers have even provided a low-level API for us to use, so we don't have to write any bit- twiddling I/O driver code.

### The hardware API

The code for the hardware interface functions written by hardware engineers can be found at:
  [Hardware API](https://github.com/anuchandy/coffeemaker/tree/master/hardwareAPI)

In the original problem statement all the hardware interface functions are described in a single interface CoffeeMakerAPI.cs.
After going through [Mark Seemann] (https://twitter.com/ploeh)'s [Encapsulation and SOLID] (http://beta.pluralsight.com/courses/encapsulation-solid)
course, I decided to apply [Command-Query seperation principle](https://en.wikipedia.org/wiki/Command–query_separation) which result in two interfaces:
 
* [CommandAPI](https://github.com/anuchandy/coffeemaker/blob/master/hardwareAPI/commandAPI.go)
* [QueryAPI] (https://github.com/anuchandy/coffeemaker/blob/master/hardwareAPI/queryAPI.go)

The HardwareAPI interface composes above two interfaces.

```go
type HardwareAPI interface {
  QueryAPI
  CommandAPI
}
```

I would highly recommand watching Mark Seemann's 'Encapsulation and SOLID' course, he uses the same Coffee Marker problem to show-case how to apply SOILD principles.

### Design

![Alt text](/CoffeeMaker.JPG?raw=true "Coffee-Maker design")

#### The Coffee machine simulator [HardwareAPIImpl & main.go]

A command-line project that simulates Coffee-Machine working can be found [here] (https://github.com/anuchandy/coffeemakerSimulator). The simulator project contains two components:

* A mock [implementation](https://github.com/anuchandy/coffeemakerSimulator/tree/master/hardwareAPIImpl) HardwareAPIImpl of above mentoned hardware API
* A [command-line interface](https://github.com/anuchandy/coffeemakerSimulator/blob/master/main.go) that allows user to:
    1. Fill water into Boiler
    2. Place pot in Warmer plate
    3. Remove pot from Warmer plate
    4. Press Brewing button
    5. See the current state of the machine
    6. Exit the simulation

The simulator creates the mock hardware, switch-on the coffee-maker [The is means coffee-maker will start monitoring and controlling various components of the hardware via API], read input from user and invokes various hardware methods to trigger user action.

```go
  // Creates a mock coffee-machine hardware
  var cmHardware *hardwareAPIImpl.HardwareAPIImpl = &hardwareAPIImpl.HardwareAPIImpl{}
  // Initializes the hardware
  cmHardware.Reset()

  // Switch-on the coffee-maker.
  coffeemaker.SwitchOn(cmHardware)

  var ui hardwareAPIImpl.UserAction = cmHardware
  for ;; {
    var action int
    fmt.Print("\nAction [1: Fill_Water 2: Place_Pot 3: Remove_Pot 4: Press_BrewButton 5: Show_Status 6: Exit] : ")
    fmt.Scanf("%d", &action)

    if action == 1 {
      ui.FillWater()
    }	else if action == 2 {
      ui.PutPot()
    } else if action == 3 {
      ui.RemovePot()
    } else if action == 4 {
      ui.PressBrewButton()
    } else if action == 5 {
      ui.ShowState()
    } else if action == 6 {
      break
    } else  {
      fmt.Println("Unknown action")
    }
  }

  // Switch-off the coffee-maker.
  coffeemaker.SwitchOff()
```

The HardwareAPIImpl satisfies an interface [UserAction](https://github.com/anuchandy/coffeemakerSimulator/blob/master/hardwareAPIImpl/hardwareAPIlmpl.go#L13), inaddition to HardwareAPI interface.

```go
type UserAction interface {
  // Fill water to boiler.
  FillWater()
  // Press the brewing button.
  PressBrewButton()
  // Place pot in the warmer plate.
  PutPot()
  // Remove pot from the warmer plate.
  RemovePot()
  // Show the current status, state of coffee-maker.
  ShowState()
}
```

#### EventAggregator

Coffee maker uses a basic implementation of EventAggregator, which act as a container for coffee machine hardware events that decouples hardware polling Go routine (the publisher) and the hardware controllers (subscribers). The implementation of EventAggregator can be found in [events](https://github.com/anuchandy/coffeemaker/tree/master/events) directory.

The published events will be send to eventsChan (```go eventsChan chan Event```) channel. Aggregator runs a dedicated Go routine to receive these events and dispatch them to subscribers.

```go
func (a *Aggregator) Start() {
  go func() {
    for {
      e, ok := <-a.eventsChan
      if !ok {
        return
      }

      for _, subscriber := range a.subscribers[e] {
        go subscriber.HandleEvent(e)
      }
    }
  }()
}
```

This routine returns when it sees the eventChan is closed. The eventChan will be closed when we switch-off the coffee-maker.

#### Coffee-machine Hardware monitoring (polling) [Publisher]

Coffee maker uses a seperate Go routine to poll the state of various hardware components and publish those states via EventAggregator.

The Go routine polls and publish events in every one second in an infinte for loop.

```go
func pollHardware(api hardwareAPI.QueryAPI) {
  ticker := time.NewTicker(1 * time.Second)
  go func() {
    for {
      select {
        case <-ticker.C:
          publishEvents(api)
        case <-abortPoll:
          ticker.Stop()
          return
      }
    }
  }()
}

// publishEvents publishes the current state of coffee-maker.
func publishEvents(api hardwareAPI.QueryAPI) {
  agg.Publish(toBoilerEvent(api.GetBoilerStatus()))
  agg.Publish(toBrewButtonEvent(api.GetBrewButtonStatus()))
  agg.Publish(toWarmerPlateEvent(api.GetWarmerPlateStatus()))
}
```

This Go routine uses [select](https://gobyexample.com/select) to listen on two channels timer and abortPoll. The routine returns When it sees an event in abortPoll channel which stops the polling.

We sent abort signal when we switch-off the coffee machine.

```go
func SwitchOff() {
  abortPoll <- struct{}{} // Stop hardware polling
  agg.Stop() // Stop aggregator from listening and publishing events
}
```

#### Coffee-machine Hardware controllers [Subscribers]

Each hardware component of the coffee-maker is controlled by seperate controllers. The controllers can be found under [hardwareController](https://github.com/anuchandy/coffeemaker/tree/master/hardwareController) directory.

Each controller subscribe for one or more hardware events, these events are used to decide:
* Whether or not to change the state of the component it controlling.
* If state needs to be changed then what should to be the new state.

While creating instance of controllers, reference to the EventAggregator and hardware API interface will be passed to the creater method.

e.g.
```go
func NewBoilerController(aggregator *events.Aggregator, api hardwareAPI.CommandAPI) *BoilerController
```

The creator method create the controller and register the controller as subscribers for the events that they want to receive.




