package logger

// func WithLogClientID(ctx context.Context, clientID string) context.Context {
// 	if c, ok := ctx.Value(key).(logCtx); ok {
// 		c.clientID = clientID
// 		return context.WithValue(ctx, key, c)
// 	}
// 	return context.WithValue(ctx, key, logCtx{clientID: clientID})
// }

// func WithLogTokenAmount(ctx context.Context, amount float64) context.Context {
// 	if c, ok := ctx.Value(key).(logCtx); ok {
// 		c.tokensAmount = amount
// 		return context.WithValue(ctx, key, c)
// 	}
// 	return context.WithValue(ctx, key, logCtx{tokensAmount: amount})
// }

// func WithLogBackendAddress(ctx context.Context, address string) context.Context {
// 	if c, ok := ctx.Value(key).(logCtx); ok {
// 		c.backendAddress = address
// 		return context.WithValue(ctx, key, c)
// 	}
// 	return context.WithValue(ctx, key, logCtx{backendAddress: address})
// }

// func WithLogFullRefill(ctx context.Context) context.Context {
// 	if c, ok := ctx.Value(key).(logCtx); ok {
// 		c.fullRefill = true
// 		return context.WithValue(ctx, key, c)
// 	}
// 	return context.WithValue(ctx, key, logCtx{fullRefill: true})

// }

// func WithLogCheck(
// 	ctx context.Context,
// 	healthUrl string,
// 	healthMethod string,
// 	err error) context.Context {
// 	if c, ok := ctx.Value(key).(logCtx); ok {
// 		c.healthUrl = healthUrl
// 		c.healthMethod = healthMethod
// 		c.err = err
// 		return context.WithValue(ctx, key, c)
// 	}
// 	return context.WithValue(ctx, key,
// 		logCtx{
// 			healthUrl:     healthUrl,
// 			healthMethod:  healthMethod,
// 			err: err,
// 		},
// 	)
// }
