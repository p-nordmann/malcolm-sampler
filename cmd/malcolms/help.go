/*
Malcolms provides a gRPC API that exposes Malcolm's Sampler functionalities.

IT is designed to meet the following requirements:
  - allow to register bounding boxes and true posterior samples of variable sizes and dimensions.
  - allow to generate samples from them in a concurrent fashion.

See  proto definition in package github.com/p-nordmann/malcolm-sampler/grpc for a more detailed
specification of the API.
*/
package main
