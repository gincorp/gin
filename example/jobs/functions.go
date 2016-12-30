package main

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "net/http"
    "net/mail"
    "net/smtp"
    "os"

    "github.com/gincorp/gin/taskmanager"
)

type financeitem struct {
    Name, Price string
}

func apiCall(url string) (d map[string]interface{}, err error) {
    var resp *http.Response

    if resp, err = http.Get(url); err != nil {
        return
    }

    buf := new(bytes.Buffer)
    buf.ReadFrom(resp.Body)
    defer resp.Body.Close()

    d = make(map[string]interface{})
    json.Unmarshal(buf.Bytes(), &d)

    return
}

func getCurrencyPrices(jn taskmanager.JobNotification) (output map[string]interface{}, err error) {
    respData, err := apiCall("http://finance.yahoo.com/webservice/v1/symbols/allcurrencies/quote?format=json")
    if err != nil {
        return
    }

    output = make(map[string]interface{})
    prices := []financeitem{}

    // Somewhere, and I suspect in java world, there's a fucking awful library turning xml into json
    // by trying to map xml classes/types to a json object. It means stupid shit like this format.
    //
    // Seriously: fuck anybody who does this.

    data := respData["list"].(map[string]interface{}) // ITS THE ONLY FUCKING KEY IN THE FUCKING OBJECT

    for _, resource := range data["resources"].([]interface{}) {
        actualResource := resource.(map[string]interface{})["resource"].(map[string]interface{}) // WANKERS
        fields := actualResource["fields"].(map[string]interface{})

        switch fields["name"].(string) {
        case "USD/GBP":
            prices = append(prices, financeitem{"$/£", fields["price"].(string)})
        case "USD/EUR":
            prices = append(prices, financeitem{"$/€", fields["price"].(string)})
        case "GOLD 1 OZ":
            prices = append(prices, financeitem{"Gold per ounce", fields["price"].(string)})
        }
    }
    output["prices"] = prices

    return
}

func sendEmail(jn taskmanager.JobNotification) (output map[string]interface{}, err error) {
    // Hat tip: https://gist.github.com/andelf/5004821
    username := os.Getenv("MAIL_USERNAME")
    password := os.Getenv("MAIL_PASSWORD")
    mailhost := jn.Context["host"]
    mailport := jn.Context["port"]

    auth := smtp.PlainAuth("",
        username,
        password,
        mailhost,
    )

    from := mail.Address{"gin", jn.Context["from"]}
    to := mail.Address{"", jn.Context["to"]}

    header := make(map[string]string)

    header["From"] = from.String()
    header["To"] = to.String()
    header["Subject"] = jn.Context["subject"]
    header["MIME-Version"] = "1.0"
    header["Content-Type"] = "text/plain; charset=\"utf-8\""
    header["Content-Transfer-Encoding"] = "base64"

    message := ""
    for k, v := range header {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(jn.Context["body"]))

    err = smtp.SendMail(
        mailhost+":"+mailport,
        auth,
        from.Address,
        []string{to.Address},
        []byte(message),
    )

    output = make(map[string]interface{})
    return
}
