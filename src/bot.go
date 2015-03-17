package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "github.com/comstud/slopher"
)

const AUTH_TOKEN = "<WHATEVER>"

type BotStateManager struct {
            *slopher.DefaultStateManager
    // Add more things here if needed
}

func (self *BotStateManager) RTMStart(rtm *slopher.RTMProcessor, resp *slopher.RTMStartResponse) {
    self.DefaultStateManager.RTMStart(rtm, resp)
    // Store more things if needed
}

func getLogger(fname string) (logger *log.Logger, err error) {
    f, err := os.OpenFile(fname,
        os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        return
    }

    logger = log.New(f, "", log.Ldate|log.Ltime|log.Lmicroseconds)
    return
}


func onMessage(rtm *slopher.RTMProcessor, _msg slopher.RTMMessage) {
    smgr := rtm.StateManager.(*BotStateManager)
    msg := _msg.(*slopher.RTMChannelMessage)

    from := smgr.FindEntity(msg.UserID)
    if from == nil {
        fmt.Printf("Message from unknown user: %s\n", msg.UserID)
        return
    }

    to := smgr.FindPlace(msg.ChannelID)
    if to == nil {
        fmt.Printf("Message to unknown channel: %s\n", msg.ChannelID)
        return
    }

    if from.IsSelf() {
        fmt.Printf("Ignoring message from self\n")
        return
    }

    if !to.IsIM {
        fmt.Printf("Ignoring message to channel %s: %s\n",
            to.Name, msg.Text)
        return
    }

    fmt.Printf("Message from %s: %s\n", from.Name, msg.Text)
}

func onTyping(rtm *slopher.RTMProcessor, _msg slopher.RTMMessage) {
    smgr := rtm.StateManager.(*BotStateManager)
    msg := _msg.(*slopher.RTMTypingMessage)

    from := smgr.FindEntity(msg.UserID)
    to := smgr.FindPlace(msg.ChannelID)

    if from.Name == to.Name {
        fmt.Printf("%s is typing a DM...\n", from.Name)
    } else {
        fmt.Printf("%s is typing to %s\n", from.Name, to.Name)
    }
}

func onChannelJoined(rtm *slopher.RTMProcessor, _msg slopher.RTMMessage) {
    smgr := rtm.StateManager.(*BotStateManager)
    msg := _msg.(*slopher.RTMChannelJoinedMessage)
    place := smgr.FindPlace(msg.Channel.ID)

    fmt.Printf("Yo, I joined %s\n", place.Name)
}

func main() {
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage: %s <logfile>\n", os.Args[0])
        flag.PrintDefaults()
        os.Exit(1)
    }
    flag.Parse()
    args := flag.Args()
    if len(args) < 1 {
        fmt.Fprintln(os.Stderr, "Error: Missing logfile")
        flag.Usage()
    }

    fname := args[0]

    logger, err := getLogger(fname)
    if err != nil {
        fmt.Printf("Couldn't get logger: %s\n", err)
        os.Exit(2)
    }

    cli := slopher.NewClient("", AUTH_TOKEN, logger)

    state_mgr := &BotStateManager{
        DefaultStateManager: slopher.GetDefaultStateManager(),
    }

    rtm_config := &slopher.RTMConfig{
        StateManager:  state_mgr,
        AutoReconnect: true,
    }

    rtm_processor, err := cli.NewRTMProcessor(rtm_config)
    if err != nil {
        fmt.Printf("Error creating RTM Processor: %s\n", err)
        os.Exit(2)
    }

    rtm_processor.OnTyping(onTyping)
    rtm_processor.OnChannelMessage(onMessage)

    if err := rtm_processor.Start(); err != nil {
        fmt.Printf("Error starting RTM processor: %s\n", err)
        os.Exit(2)
    }

    <-rtm_processor.GetDoneChannel()
}
