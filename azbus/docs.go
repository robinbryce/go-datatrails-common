package azbus

// TODO: this needs some attention as it is incorrect

// Implement pubsubinterface
// Azure service bus pubsub interface
//
// Abstracts away all the Azure plumbing and replaces with simpler interface.
//
// First create the interface:
//
//
// Usage: Sending message(s):
//	client := azsb.Receiver{
//		Cfg: azsb.ReceiverConfig{
//			ConnectionString: "blah-blah-blah...",
//			TopicName: "userinterface",
//		},
//	}
//
//  ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
//  defer cancel()
//  client.Open(ctx)
//  defer client.Close(ctx)
//  client.Publish(ctx, []byte("Hello World"))
//  ...... more send messages if required.
//
//	Receiving messages:
//
//	client := azsb.Receiver{
//		Cfg: azsb.ReceiverConfig{
//			ConnectionString: "blah-blah-blah...",
//			TopicName: "userinterface",
//			SubscriptionName: "xxxx",
//		},
//	}

//
//  handler := func(ctx context.Context, msg *azbus.ReceivedMessage) error {
//		logger.Sugar.Infow(
//			"Received",
//			"message", msg,
//		)
//      ctx, cancel := contexttWithTimeout(ctx, 60*time.Second)
//      defer cancel()
//      // do stuff
//		msg.Complete(ctx)
//		return nil
// 	}
//  err := client.Subscribe(ctx, handler)
//  var azerr *azbus.AzbusError
//  if errors.As(err, &azerr) {
//      // TopicName, SubscriptionName may be wrong so die...
//      // else service bus not yet available so die...
//  } else {
//      // log transient error - maybe sleep and retry..
//  }
