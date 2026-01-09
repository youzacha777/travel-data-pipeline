
```
project05
└─ travel-data-pipeline
   ├─ docker
   │  ├─ clickhouse-init
   │  │  └─ create-tables.sql
   │  ├─ docker-compose.yml
   │  ├─ Dockerfile.flink
   │  ├─ Dockerfile.kafka
   │  ├─ flink-checkpoints
   │  │  ├─ 0c7b8b3ce1ab2afa01baf6518d2fdb5f
   │  │  │  ├─ shared
   │  │  │  └─ taskowned
   │  ├─ flink-init
   │  │  └─ flink-job.sh
   │  ├─ flink-kafka-clickhouse
   │  │  └─ target
   │  │     └─ flink-job-init.log
   │  ├─ kafka-init
   │  │  ├─ create-topics.sh
   │  │  └─ format-and-start.sh
   │  └─ libs
   │     ├─ clickhouse-jdbc-0.6.3-all.jar
   │     ├─ flink-connector-jdbc-3.3.0-1.20.jar
   │     ├─ flink-json-1.20.1.jar
   │     └─ flink-sql-connector-kafka-3.3.0-1.20.jar
   ├─ event-generator
   │  ├─ cmd
   │  │  └─ generator
   │  │     └─ main.go
   │  ├─ go.mod
   │  ├─ go.sum
   │  └─ internal
   │     ├─ controller
   │     │  └─ load_controller.go
   │     ├─ event
   │     │  └─ event.go
   │     ├─ fsm
   │     │  ├─ fsm.go
   │     │  ├─ transitions.go
   │     │  ├─ types.go
   │     │  └─ utils.go
   │     ├─ generator
   │     │  ├─ addtocart.go
   │     │  ├─ browsing.go
   │     │  ├─ click.go
   │     │  ├─ eventbrowsing.go
   │     │  ├─ nextpage.go
   │     │  ├─ payload.go
   │     │  ├─ product_catalog.go
   │     │  ├─ purchase.go
   │     │  ├─ search.go
   │     │  └─ utils.go
   │     ├─ metrics
   │     │  ├─ inmemory.go
   │     │  └─ metrics.go
   │     ├─ user
   │     │  ├─ session.go
   │     │  ├─ session_manager.go
   │     │  └─ user_pool.go
   │     └─ worker
   │        └─ worker.go
   └─ flink-kafka-clickhouse
      ├─ libs
      │  └─ clickhouse-jdbc-0.6.3-all.jar
      ├─ pom.xml
      ├─ src
      │  └─ main
      │     └─ java
      │        └─ com
      │           └─ flink
      │              └─ FlinkKafkaToClickhouse.java
      └─ target
         ├─ classes
         │  └─ com
         │     └─ flink
         │        └─ FlinkKafkaToClickhouse.class
         ├─ flink-kafka-clickhouse-1.0-SNAPSHOT.jar
         ├─ generated-sources
         │  └─ annotations
         ├─ maven-archiver
         │  └─ pom.properties
         ├─ maven-status
         │  └─ maven-compiler-plugin
         │     └─ compile
         │        └─ default-compile
         │           ├─ createdFiles.lst
         │           └─ inputFiles.lst
         ├─ original-flink-kafka-clickhouse-1.0-SNAPSHOT.jar
         └─ test-classes

```