// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
// SPDX-FileCopyrightText: 2019-present Rimedo Labs
//
// SPDX-License-Identifier: Apache-2.0

package utils

const (
	// MhoServiceModelName is MHO service model name
	MhoServiceModelName = "oran-e2sm-mho"
	// MhoServiceModelVersion is the version of MHO SM
	MhoServiceModelVersion = "v2"

	// TLSCacrt is TLS cacrt value
	TLSCacrt = "-----BEGIN CERTIFICATE-----\nMIIDYDCCAkgCCQDe99fSN9qxSTANBgkqhkiG9w0BAQsFADByMQswCQYDVQQGEwJV\nUzELMAkGA1UECAwCQ0ExEjAQBgNVBAcMCU1lbmxvUGFyazEMMAoGA1UECgwDT05G\nMRQwEgYDVQQLDAtFbmdpbmVlcmluZzEeMBwGA1UEAwwVY2Eub3Blbm5ldHdvcmtp\nbmcub3JnMB4XDTE5MDQxMTA5MDYxM1oXDTI5MDQwODA5MDYxM1owcjELMAkGA1UE\nBhMCVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlNZW5sb1BhcmsxDDAKBgNVBAoM\nA09ORjEUMBIGA1UECwwLRW5naW5lZXJpbmcxHjAcBgNVBAMMFWNhLm9wZW5uZXR3\nb3JraW5nLm9yZzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMEg7CZR\nX8Y+syKHaQCh6mNIL1D065trwX8RnuKM2kBwSu034zefQAPloWugSoJgJnf5fe0j\nnUD8gN3Sm8XRhCkvf67pzfabgw4n8eJmHScyL/ugyExB6Kahwzn37bt3oT3gSqhr\n6PUznWJ8fvfVuCHZZkv/HPRp4eyAcGzbJ4TuB0go4s6VE0WU5OCxCSlAiK3lvpVr\n3DOLdYLVoCa5q8Ctl3wXDrfTLw5/Bpfrg9fF9ED2/YKIdV8KZ2ki/gwEOQqWcKp8\n0LkTlfOWsdGjp4opPuPT7njMBGXMJzJ8/J1e1aJvIsoB7n8XrfvkNiWL5U3fM4N7\nUZN9jfcl7ULmm7cCAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAIh6FjkQuTfXddmZY\nFYpoTen/VD5iu2Xxc1TexwmKeH+YtaKp1Zk8PTgbCtMEwEiyslfeHTMtODfnpUIk\nDwvtB4W0PAnreRsqh9MBzdU6YZmzGyZ92vSUB3yukkHaYzyjeKM0AwgVl9yRNEZw\nY/OM070hJXXzJh3eJpLl9dlUbMKzaoAh2bZx6y3ZJIZFs/zrpGfg4lvBAvfO/59i\nmxJ9bQBSN3U2Hwp6ioOQzP0LpllfXtx9N5LanWpB0cu/HN9vAgtp3kRTBZD0M1XI\nCtit8bXV7Mz+1iGqoyUhfCYcCSjuWTgAxzir+hrdn7uO67Hv4ndCoSj4SQaGka3W\neEfVeA==\n-----END CERTIFICATE-----"
	// TLSKey is TLS key value
	TLSKey = "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC42bpv4K93I1QN\nSGpTxU26V0BUm2+1ciebv5YIKOlcm1gHD7YCaxrVLbRjmHoeccjFcSyYHzB6FD2O\nXTcJTxndiBcjswk2MUQAdLsAo18rPtJL3bGgDjrC6RXeeOY3jGT8ll8WttpQngdw\nTk9WLl0nLJDGDqY9yl1uJgdDj64M3hwpTyYJFZ2PiA1RoL6bz6oeq3dXdV9ryIEx\n2myRPBP7WMFBwMa34DLaKKcICZ/mbSmASRpWCG1cgotSMb9rXsUCSDL+xNYnsnjB\noPSJFCwQF5D7ivYwIi20tGT+ZNN0Es6QZSK6fON6H7USIWyU1TBfG7Sp+GEkHco1\n/B7RsX4tAgMBAAECggEBAIvky0HsKx7g77V1vnJTebWyXo8pa2tIT02BusvGGoXp\nUr9VVouR/yaihkhxlsn/ltBGDFe8EvXw530ccpBq+so7OjfcQPZwZmRp8zRSb63M\nx15/EvRskG/98n0Bxkj3yV2Xd7M7AxHL5xlJSqWQRRNmmNIrOAi/Y+H+ibTJwhEd\npV6pULHmvu4mU2YdkX/6RLOS89aAxOzXgs6nUzp826/ugb4X34izp/WwMzCs+QJd\nQVvjA/zuLnzIqspqt9bNJ3bCs+/ovqpdYlDBYbYHA0kLEMYOQsAkrWjJWPoLoDyD\nHLZZgG9jkQzkxdUzquku4UahsZZi8317jyT5b/GX0AECgYEA43jvJk887exwc04q\ns2/ef0spPD4JGVmblyD1upWdiWNgyRxxNbOkqaR1Gp5TqcdCmNRaraP3ZKbpngmp\nQKvKGLf/RH7H1b9/NuAlLAI8rSiP6h1MDyTbO1wwdSU4f1HO/3pVyRJfXd2h8+eD\nVfXe8rfXafWBBXWOeJ52JuDStC0CgYEA0Ahn7ikHKo8HM2FIJkrhm2Y4jwSL9TvO\nS6ikBbQPUbhcdyLbHVrOLlmiiBqNnuqe8AzaHrH+T6TeP312M5GGVT9CUkIb1Vq6\nfC8yVTDg4gPlSuz974xRSujWWLutX/6Hr8eAN8L/E78LD/Ojy+DISnIzmTC2B5lM\no6WJ+BcmMgECgYAxrrY9Hc1nAd9Fr+rvqh1knBvzhnEiUkoDZjWFfSwdV9FJ26Z2\nXjg2vS6+k5oeWOEY1DjB+DAOkc4wsFeBQoQvhfCBG1e2Pc8hQy+bPxnVkChur9tu\n61Pe0THcRDbkyA94CVY3RoYB0GiRBx3OZpc9WB36jJ6TfKuTeLjBoRUkOQKBgHl9\nPzy9pxq6lojx+hGqz2BSbRtQm2+m8o4KuWc/RWcDFLTanT3iZuB4pkt3vlcdS56C\n0ur0JcFbVhOb8GijRuEH5XJmexy5NIkLgwhvWBWGEuUTzCSWPG9T1MHTMKgL3C/S\ngVWPQinE+u/g6DpLVozrbqi64sNDSpeTOCSzWDIBAoGAWBxIN1k+xQQYneI1+uvi\nd8NpUcLRilayIKQkaZFA+efXpyPOq8r0WtAh8tpTA3NFXunMiejJUF1wkzL4+NLt\n3ct4RY4eHoPQHtxOqZn7aMx+8/V4yz6IDKEsAsxK1p7AP6yt6GEWzj8OvrP7cm4Z\nNVQirzdY6fbkOULBISdVSWk=\n-----END PRIVATE KEY-----"
	//TLSCrt is TLS crt value
	TLSCrt = "-----BEGIN CERTIFICATE-----\nMIIDcTCCAlkCFErBGzsXHo1l8bmZRmDkF+h2bsdVMA0GCSqGSIb3DQEBCwUAMHIx\nCzAJBgNVBAYTAlVTMQswCQYDVQQIDAJDQTESMBAGA1UEBwwJTWVubG9QYXJrMQww\nCgYDVQQKDANPTkYxFDASBgNVBAsMC0VuZ2luZWVyaW5nMR4wHAYDVQQDDBVjYS5v\ncGVubmV0d29ya2luZy5vcmcwHhcNMjAwOTAxMDYwNTI2WhcNMzAwODMwMDYwNTI2\nWjB4MQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExEjAQBgNVBAcMCU1lbmxvUGFy\nazEMMAoGA1UECgwDT05GMRQwEgYDVQQLDAtFbmdpbmVlcmluZzEkMCIGA1UEAwwb\nb25vcy1lMnQub3Blbm5ldHdvcmtpbmcub3JnMIIBIjANBgkqhkiG9w0BAQEFAAOC\nAQ8AMIIBCgKCAQEAuNm6b+CvdyNUDUhqU8VNuldAVJtvtXInm7+WCCjpXJtYBw+2\nAmsa1S20Y5h6HnHIxXEsmB8wehQ9jl03CU8Z3YgXI7MJNjFEAHS7AKNfKz7SS92x\noA46wukV3njmN4xk/JZfFrbaUJ4HcE5PVi5dJyyQxg6mPcpdbiYHQ4+uDN4cKU8m\nCRWdj4gNUaC+m8+qHqt3V3Vfa8iBMdpskTwT+1jBQcDGt+Ay2iinCAmf5m0pgEka\nVghtXIKLUjG/a17FAkgy/sTWJ7J4waD0iRQsEBeQ+4r2MCIttLRk/mTTdBLOkGUi\nunzjeh+1EiFslNUwXxu0qfhhJB3KNfwe0bF+LQIDAQABMA0GCSqGSIb3DQEBCwUA\nA4IBAQCAron90Id5rzqn73M7FcCN9pFtu1MZ/WYNBxPmrcRc/yZ80PecZoHgTnJh\nmBDTLwpoRLPimxTL4OzrnA6Go0kD/CPAThehGb8BBZ+aiSJ17I0/EL1HDmXStgRk\nWuqP2DxenckWHaNmPVE0PbB6BCsd5HP0tCC4vGBbGYbJmAhhjzhzEmEypqskt+Np\neFe1DgDyfVrroIHmDLPCEu2ny9Syr/LslDmndGses8/QSVDDyAK/LFFMukCJRWsQ\nuIUJM/aDEAqbZUs4bb60hVfcZTU1HVPcp2xuOmVUFKUvHyCpt/n65Y/5XKQYQpTr\n1qa1krCQOnuwSstIpqBCnX+TecP7\n-----END CERTIFICATE-----"
	//ConfigJSON has JSON-type config parameters
	ConfigJSON = "{\n	\"reportingPeriod\": 1000,\n	\"periodic\": true,\n	\"uponRcvMeasReport\": true,\n	\"uponChangeRrcStatus\": true,\n	\"A3OffsetRange\": 0,\n	\"HysteresisRange\": 0,\n	\"CellIndividualOffset\": 0,\n	\"FrequencyOffset\": 0,\n	\"TimeToTrigger\": 0\n}"

	// Verification time
	VerificationTimer = 120 // (Seconds)

	// Serving cell
	SCellID = "138426014550001"

	A1Policy = `{
	   "scope":{
		  "ueId":"<IMSI>"
	   },
	   "tspResources":[
		  {
			 "cellIdList":[
				{
				   "plmnId":{
					  "mcc":"138",
					  "mnc":"426"
				   },
				   "cId":{
					  "ncI":470106432
				   }
				}
			 ],
			 "preference":"FORBID"
		  }
	   ]
	}`
)
