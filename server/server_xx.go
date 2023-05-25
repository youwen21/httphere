package server

//func getRewriteRule(cfg interface{}) map[string]string {
//	confMap, ok := cfg.(map[string]interface{})
//	if !ok {
//		return nil
//	}
//	ruleI, ok := confMap["rewrite"]
//	if !ok {
//		return nil
//	}
//	rule, ok := ruleI.(map[string]interface{})
//	if !ok {
//		return nil
//	}
//
//	result := make(map[string]string)
//	for k, v := range rule {
//		result[k] = cast.ToString(v)
//	}
//	return result
//}

//func getPaths(cfg interface{}) map[string]string {
//	confMap, ok := cfg.(map[string]interface{})
//	if !ok {
//		return nil
//	}
//
//	result := make(map[string]string)
//	for k, v := range confMap {
//		if k == "rewrite" {
//			continue
//		}
//		if vStr, vok := v.(string); vok {
//			result[k] = vStr
//		}
//	}
//
//	return result
//}

//func RewritePath(path string, rewriteRule map[string]string) string {
//	for k, v := range rewriteRule {
//		if strings.HasPrefix(path, k) {
//			return strings.Replace(path, k, v, 1)
//		}
//	}
//
//	return path
//}

//func (f MyServer) RewritePath(path string) string {
//	if f.rewrite == nil {
//		return path
//	}
//
//	for k, v := range f.rewrite {
//		if strings.HasPrefix(path, k) {
//			return strings.Replace(path, k, v, 1)
//		}
//	}
//
//	return path
//}

//func initMuxServerByConf(hostConf interface{}) *http.ServeMux {
//	pathsCfg := getPaths(hostConf)
//	if pathsCfg == nil {
//		return nil
//	}
//	server := http.NewServeMux()
//	for k, v := range pathsCfg {
//		backendURL, err := url.Parse(v)
//		if err != nil {
//			fmt.Println("config has error:", hostConf)
//			continue
//		}
//		server.Handle(k, NewSingleHostReverseProxyFake(backendURL))
//	}
//
//	return server
//}
