Distribyted.logs = {
    loadView: function () {
        fetch("/api/log")
            .then(response => {
                if (response.ok) {
                    return response.body.getReader();
                } else {
                    response.json().then(json => {
                        Distribyted.message.error('Error getting logs from server. Error: ' + json.error);
                    }).catch(error => {
                        Distribyted.message.error('Error getting logs from server. Error: ' + error);
                    })
                }
            })
            .then(reader => {
                var decoder = new TextDecoder()
                var lastString = ''
                reader.read().then(function processText({ done, value }) {
                    if (done) {
                        return;
                    }

                    const string = `${lastString}${decoder.decode(value)}`
                    const lines = string.split(/\r\n|[\r\n]/g)
                    this.lastString = lines.pop() || ''

                    lines.forEach(element => {
                        try {
                            var json = JSON.parse(element)
                            var properties = ""
                            for (let [key, value] of Object.entries(json)) {
                                if (key == "level" || key == "component" || key == "message" || key == "time") {
                                    continue
                                }

                                properties += `<b>${key}</b>=${value} `
                            }

                            var tableClass = "table-primary"
                            switch (json.level) {
                                case "info":
                                    tableClass = ""
                                    break;
                                case "error":
                                    tableClass = "table-danger"
                                    break;
                                case "warn":
                                    tableClass = "table-warning"
                                    break;
                                case "debug":
                                    tableClass = "table-info"
                                    break;
                                default:
                                    break;
                            }
                            template = `<tr class="${tableClass}"><td>${new Date(json.time*1000).toLocaleString()}</td><td>${json.level}</td><td>${json.component}</td><td>${json.message}</td><td>${properties}</td></tr>`;
                            document.getElementById("log_table").innerHTML += template;
                        } catch (err) {
                            // server can send some corrupted json line
                            console.log(err);
                        }
                    });


                    return reader.read().then(processText);
                }).catch(err => console.log(err));
            }).catch(err => console.log(err));
    }
}
