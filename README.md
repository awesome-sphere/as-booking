# AS Booking micro service

	Service for updating status to seat database.

### How to run service
run the following command \
```
go run main.go
```
**Port:** 9009

### Required deployment variable
<ul>
	<li> KAFKA_BOOKING_TOPIC_NAME: name for booking topic</li>
	<li> KAFKA_CANCEL_TOPIC_NAME: name for cancel booking topic</li>
	<li> KAFKA_ADDR: Kafka brokers url</li>
	<li> KAFKA_BOOKING_TOPIC_PARTITION: numbers of partitions for booking topic</li>
	<li> KAFKA_CANCEL_TOPIC_PARTITION: numbers of partitions for booking topic</li>
	<li>POSTGRES_USER: database username</li>
	<li>POSTGRES_PASSWORD: database password</li>
	<li>POSTGRES_HOST: database host i.e. localhost</li>
	<li>POSTGRES_DB: database name</li>
	<li>THEATER_AMOUNT: numbers of theaters</li>

</ul>

