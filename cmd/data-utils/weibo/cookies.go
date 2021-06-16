package weibo

type cookieManager struct {
}

func NewCookieManager() *cookieManager {
	return &cookieManager{}
}

func (t *cookieManager) GetCookies(domain string) string {
	return "_s_tentry=login.sina.com.cn; Apache=8992117274353.781.1620035119646; SINAGLOBAL=8992117274353.781.1620035119646; ULV=1620035119761:1:1:1:8992117274353.781.1620035119646:; login_sid_t=ec6557efaff3fad8071a2ac25775118b; cross_origin_proto=SSL; appkey=; WBtopGlobal_register_version=91c79ed46b5606b9; YF-V-WEIBO-G0=35846f552801987f8c1e8f7cec0e2230; XSRF-TOKEN=ER7YQ1bLjtUilo5KfEPWnuf3; SSOLoginState=1623140936; UOR=login.sina.com.cn,s.weibo.com,tophub.today; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9Wh17wQNlYaBrxevjh.g7L0_5JpX5KMhUgL.FoMcehz4eh-7So-2dJLoI7ybwgDbTKqcSntt; ALF=1655369134; SCF=AnE8yPJm9c3c7KMKhrXUWPLW3K9IAAhAeQclk2MY3uPukhzc1DkcZ4-26qUIAM78fFMGvP2LuyqQNciMdUFOG9c.; SUB=_2A25Nzcp_DeRhGeFI61AY8CvMzTmIHXVuury3rDV8PUNbmtAKLRnikW9NfWSjH5_10aG5qbgOCGIdgbfGvOflGBOL; WBPSESS=em3Wpwxtl1qzqY8N-7jSOBWKDInZaOLeTW4EGUbwOAqSF0jHU_Vq69eJNmqzuqlofTy1ErF2e39GZyk8-q74CGVBYmGq1gbtPLRuuHWJUcKFzXxAFF-4SXuArfWK5QGU"
}
