package supplier

type Option func(ISupplier)

// func SetTimeOut(ttl time.Duration) Option {
// 	return func(o *options) {
// 		o.ttl = ttl
// 	}
// }
