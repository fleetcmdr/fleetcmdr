package main

func (svc *service) bindRoutes() {
	svc.router.GET("/", svc.baseHandler)

	svc.router.GET("/api/v1/parts/leftNav", svc.leftNavHandler)
	svc.router.GET("/api/v1/parts/agents/:id", svc.viewAgentHandler)

}
