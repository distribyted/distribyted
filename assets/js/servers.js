Handlebars.registerHelper("to_date", function (timestamp) {
    return new Date(timestamp * 1000).toLocaleString()
});

Distribyted.servers = {
    _template: null,

    _getTemplate: function () {
        if (this._template != null) {
            return this._template
        }

        const tTemplate = fetch('/assets/templates/servers.html')
            .then((response) => {
                if (response.ok) {
                    return response.text();
                } else {
                    Distribyted.message.error('Error getting data from server. Response: ' + response.status);
                }
            })
            .then((t) => {
                return Handlebars.compile(t);
            })
            .catch(error => {
                Distribyted.message.error('Error getting servers template: ' + error.message);
            });

        this._template = tTemplate;
        return tTemplate;
    },

    _getRoutesJson: function () {
        return fetch('/api/servers')
            .then(function (response) {
                if (response.ok) {
                    return response.json();
                } else {
                    Distribyted.message.error('Error getting data from server. Response: ' + response.status)
                }
            })
            .catch(function (error) {
                Distribyted.message.error('Error getting status info: ' + error.message)
            });
    },

    loadView: function () {
        this._getTemplate()
            .then(t =>
                this._getRoutesJson().then(routes => {
                    document.getElementById('template_target').innerHTML = t(routes);
                })
            );
    }
}