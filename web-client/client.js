const { HelloRequest, HelloReply } = require("./helloworld_pb.js");
const { GreeterClient } = require("./helloworld_grpc_web_pb.js");

window.addEventListener("load", () => {
  // let client = new GreeterClient("https://localhost:8080");
  let client = new GreeterClient("http://localhost:8080");
  let req = new HelloRequest();
  req.setName("world");

  console.log("starting");
  var deadline = new Date();
  deadline.setSeconds(deadline.getSeconds() + 1);
  metadata = { deadline: deadline.getTime() };

  call = client.sayHello(req, metadata, (err, helloReply) => {
    if (err) {
      console.log(err);
    } else {
      console.log(helloReply.getMessage());
      // window.res = res;
    }
  });
  //
  // call.on("status", s => {
  //   console.log(s);
  // });
  //
  r = new proto.helloworld.RepeatHelloRequest();
  r.setName("repeater");
  r.setCount(10);
  console.log(r);

  stream = client.sayRepeatHello(r, {});
  stream.on("data", res => {
    console.log(res.getMessage());
  });
  stream.on("status", s => {
    console.log(s.code);
    console.log(s.details);
    console.log(s.metadata);
  });
  stream.on("end", e => {
    console.log("ended");
  });
});
