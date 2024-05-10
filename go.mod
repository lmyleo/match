module match

go 1.20

require github.com/bytedance/mockey v1.2.10

require (
	github.com/360EntSecGroup-Skylar/excelize/v2 v2.8.1 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.1 // indirect
	github.com/smarty/assertions v1.15.0 // indirect
	github.com/smartystreets/goconvey v1.8.1 // indirect
	github.com/xuri/efp v0.0.0-20220407160117-ad0f7a785be8 // indirect
	github.com/xuri/nfp v0.0.0-20220409054826-5e722a1d9e22 // indirect
	golang.org/x/arch v0.0.0-20201008161808-52c3e6f60cff // indirect
	golang.org/x/crypto v0.0.0-20220408190544-5352b0902921 // indirect
	golang.org/x/net v0.0.0-20220407224826-aac1ed45d8e3 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace (
	github.com/360EntSecGroup-Skylar/excelize/v2 => github.com/xuri/excelize/v2 v2.6.0
	github.com/xuri/excelize/v2 => github.com/360EntSecGroup-Skylar/excelize/v2 v2.6.0
)
