sub Main()
    print "StreamVerse Roku App Starting..."

    ' Show splash screen
    showChannelSGScreen()
end sub

sub showChannelSGScreen()
    screen = CreateObject("roSGScreen")
    port = CreateObject("roMessagePort")
    screen.setMessagePort(port)

    ' Create main scene
    scene = screen.CreateScene("MainScene")
    screen.show()

    ' Event loop
    while(true)
        msg = wait(0, port)
        msgType = type(msg)

        if msgType = "roSGScreenEvent" then
            if msg.isScreenClosed() then
                return
            end if
        end if
    end while
end sub
