'
' Login View Controller
'

function Init()
    m.emailKeyboard = m.top.findNode("emailKeyboard")
    m.passwordKeyboard = m.top.findNode("passwordKeyboard")
    
    m.emailKeyboard.ObserveField("buttonSelected", "OnEmailButtonSelected")
    m.passwordKeyboard.ObserveField("buttonSelected", "OnPasswordButtonSelected")
    
    m.emailKeyboard.SetFocus(true)
end function

function OnEmailButtonSelected()
    if m.emailKeyboard.buttonSelected = 0 ' OK button
        m.email = m.emailKeyboard.text
        m.emailKeyboard.visible = false
        m.passwordKeyboard.visible = true
        m.passwordKeyboard.SetFocus(true)
    end if
end function

function OnPasswordButtonSelected()
    if m.passwordKeyboard.buttonSelected = 0 ' Login button
        m.password = m.passwordKeyboard.text
        PerformLogin()
    else ' Cancel
        m.emailKeyboard.SetFocus(true)
    end if
end function

function PerformLogin()
    ' Create login request
    request = CreateObject("roUrlTransfer")
    request.SetUrl("https://api.streamverse.com/api/v1/auth/login")
    request.SetRequest("POST")
    
    ' Set headers
    request.AddHeader("Content-Type", "application/json")
    request.SetCertificatesFile("common:/certs/ca-bundle.crt")
    request.InitClientCertificates()
    
    ' Create request body
    body = FormatJson({
        email: m.email
        password: m.password
    })
    
    response = request.PostFromString(body)
    
    if response = 200
        responseData = ParseJson(request.GetToString())
        if responseData <> invalid and responseData.token <> invalid
            ' Save token
            sec = CreateObject("roRegistrySection", "StreamVerse")
            sec.Write("auth_token", responseData.token)
            sec.Flush()
            
            ' Signal login success
            m.top.loginSuccess = true
        end if
    else
        ShowErrorDialog("Login failed. Please check your credentials.")
    end if
end function

function ShowErrorDialog(message as String)
    dialog = CreateObject("roSGNode", "StandardMessageDialog")
    dialog.title = "Login Error"
    dialog.message = message
    dialog.buttons = ["OK"]
    m.top.dialog = dialog
end function

