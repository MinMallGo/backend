package service

//func SpecValueCreate(c *gin.Context, create *dto.SpecValueCreate) {
//	err := dao.NewSpecValueDao().Create(create)
//	if err != nil {
//		response.Failure(c, err.Error())
//		return
//	}
//	response.Success(c, []string{})
//	return
//}
//
//func SpecValueDelete(c *gin.Context, delete *dto.SpecValueDelete) {
//	err := dao.NewSpecValueDao().Delete(delete)
//	if err != nil {
//		response.Failure(c, err.Error())
//		return
//	}
//	response.Success(c, []string{})
//	return
//}
//
//func SpecValueUpdate(c *gin.Context, update *dto.SpecValueUpdate) {
//	err := dao.NewSpecValueDao().Update(update)
//	if err != nil {
//		response.Failure(c, err.Error())
//		return
//	}
//	response.Success(c, []string{})
//	return
//}
//
//func SpecValueSearch(c *gin.Context, search *dto.SpecValueSearch) {
//	res, err := dao.NewSpecValueDao().Paginate(search)
//	if err != nil {
//		response.Failure(c, err.Error())
//		return
//	}
//	response.Success(c, res)
//	return
//}
