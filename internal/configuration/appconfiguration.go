package configuration

type AppConfiguration struct {
	DbConfiguration     *DbConfiguration     `json:"dbConfiguration"`
	EventsConfiguration *EventsConfiguration `json:"eventsConfiguration"`
}

type EventsConfiguration struct {
	BrokerUrl string `json:"brokerUrl"`
	Topic     string `json:"topic"`
}

type DbConfiguration struct {
	ConnectionRetryCount   int    `json:"connectionRetryCount"`
	ConnectionRetryTimeout int    `json:"connectionRetryTimeout"`
	ConnectionString       string `json:"connectionString"`
}
