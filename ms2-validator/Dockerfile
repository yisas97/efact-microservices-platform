FROM maven:3.9-eclipse-temurin-17 AS build

WORKDIR /app

COPY pom.xml .
RUN mvn dependency:go-offline

COPY src ./src
RUN mvn clean package -DskipTests

FROM eclipse-temurin:17-jre-alpine

RUN apk add --no-cache curl jq

WORKDIR /app

COPY --from=build /app/target/validator-1.0.0.jar app.jar

COPY init-vault-secrets.sh /app/init-vault-secrets.sh
RUN chmod +x /app/init-vault-secrets.sh

EXPOSE 8080

ENTRYPOINT ["/app/init-vault-secrets.sh"]
