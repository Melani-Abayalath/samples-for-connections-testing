import ballerina/http;

type Greeting record {
    string 'from;
    string to;
    string message;
};

service / on new http:Listener(8090) {

    resource function get greeting() returns Greeting {
        return {
            "from": "Choreo",
            "to": "hansii",
            "message": "Welcome to Choreo v1.2!"
        };
    }

    resource function get test() returns Greeting {
        return {
            "from": "Choreo",
            "to": "hansi",
            "message": "Welcome to Choreo!"
        };
    }

    // Health endpoint
    resource function get health() returns json {
        return {
            status: "UP",
            'service: "http-api"
        };
    }
}
