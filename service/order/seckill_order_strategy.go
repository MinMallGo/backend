package order

type SeckillStrategy struct{}

func (s *SeckillStrategy) CreateOrder(ctx *Context) (res Result, err error) {
	panic("implement me")
}

func (s *SeckillStrategy) orderExpire() error {
	panic("implement me")
}
