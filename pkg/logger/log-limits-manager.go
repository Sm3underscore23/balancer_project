package logger

// func WithLogClientLimits(ctx context.Context, clientID string, capacity float64, rate float64) context.Context {
// 	if c, ok := ctx.Value(key).(logCtx); ok {
// 		c.clientID = clientID
// 		c.clientCapacityLimits = capacity
// 		c.clientRateLimits = rate
// 		return context.WithValue(ctx, key, c)
// 	}
// 	return context.WithValue(ctx, key,
// 		logCtx{
// 			clientID:             clientID,
// 			clientCapacityLimits: capacity,
// 			clientRateLimits:     rate,
// 		})
// }

// func WithLogErr(ctx context.Context, err error) context.Context {
// 	if c, ok := ctx.Value(key).(logCtx); ok {
// 		c.err = err
// 		return context.WithValue(ctx, key, c)
// 	}
// 	return context.WithValue(ctx, key, logCtx{err: err})
// }
