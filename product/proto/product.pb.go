// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: product.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Запрос для поиска товаров
type SearchProductsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Query         string                 `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`    // строка поиска
	Limit         int32                  `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`   // лимит результатов
	Offset        int32                  `protobuf:"varint,3,opt,name=offset,proto3" json:"offset,omitempty"` // смещение для пагинации
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchProductsRequest) Reset() {
	*x = SearchProductsRequest{}
	mi := &file_product_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchProductsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchProductsRequest) ProtoMessage() {}

func (x *SearchProductsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchProductsRequest.ProtoReflect.Descriptor instead.
func (*SearchProductsRequest) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{0}
}

func (x *SearchProductsRequest) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *SearchProductsRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *SearchProductsRequest) GetOffset() int32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

// Ответ по поиску товаров
type SearchProductsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Products      []*Product             `protobuf:"bytes,1,rep,name=products,proto3" json:"products,omitempty"`
	TotalCount    int32                  `protobuf:"varint,2,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"` // общее количество
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchProductsResponse) Reset() {
	*x = SearchProductsResponse{}
	mi := &file_product_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchProductsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchProductsResponse) ProtoMessage() {}

func (x *SearchProductsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchProductsResponse.ProtoReflect.Descriptor instead.
func (*SearchProductsResponse) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{1}
}

func (x *SearchProductsResponse) GetProducts() []*Product {
	if x != nil {
		return x.Products
	}
	return nil
}

func (x *SearchProductsResponse) GetTotalCount() int32 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

// Запрос для получения товара по ID
type GetProductRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProductId     int32                  `protobuf:"varint,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetProductRequest) Reset() {
	*x = GetProductRequest{}
	mi := &file_product_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetProductRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProductRequest) ProtoMessage() {}

func (x *GetProductRequest) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProductRequest.ProtoReflect.Descriptor instead.
func (*GetProductRequest) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{2}
}

func (x *GetProductRequest) GetProductId() int32 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

// Запрос для добавления нового товара
type AddProductRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Title         string                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Price         float64                `protobuf:"fixed64,3,opt,name=price,proto3" json:"price,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddProductRequest) Reset() {
	*x = AddProductRequest{}
	mi := &file_product_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddProductRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddProductRequest) ProtoMessage() {}

func (x *AddProductRequest) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddProductRequest.ProtoReflect.Descriptor instead.
func (*AddProductRequest) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{3}
}

func (x *AddProductRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *AddProductRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *AddProductRequest) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

// Ответ при добавлении товара
type AddProductResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProductId     int32                  `protobuf:"varint,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddProductResponse) Reset() {
	*x = AddProductResponse{}
	mi := &file_product_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddProductResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddProductResponse) ProtoMessage() {}

func (x *AddProductResponse) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddProductResponse.ProtoReflect.Descriptor instead.
func (*AddProductResponse) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{4}
}

func (x *AddProductResponse) GetProductId() int32 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

// Запрос для удаления товара по ID
type DeleteProductRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProductId     int32                  `protobuf:"varint,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteProductRequest) Reset() {
	*x = DeleteProductRequest{}
	mi := &file_product_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteProductRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteProductRequest) ProtoMessage() {}

func (x *DeleteProductRequest) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteProductRequest.ProtoReflect.Descriptor instead.
func (*DeleteProductRequest) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteProductRequest) GetProductId() int32 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

// Ответ при удалении товара
type DeleteProductResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteProductResponse) Reset() {
	*x = DeleteProductResponse{}
	mi := &file_product_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteProductResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteProductResponse) ProtoMessage() {}

func (x *DeleteProductResponse) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteProductResponse.ProtoReflect.Descriptor instead.
func (*DeleteProductResponse) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteProductResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

// Модель товара
type Product struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Price         float64                `protobuf:"fixed64,4,opt,name=price,proto3" json:"price,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Product) Reset() {
	*x = Product{}
	mi := &file_product_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Product) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Product) ProtoMessage() {}

func (x *Product) ProtoReflect() protoreflect.Message {
	mi := &file_product_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Product.ProtoReflect.Descriptor instead.
func (*Product) Descriptor() ([]byte, []int) {
	return file_product_proto_rawDescGZIP(), []int{7}
}

func (x *Product) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Product) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Product) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Product) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

var File_product_proto protoreflect.FileDescriptor

const file_product_proto_rawDesc = "" +
	"\n" +
	"\rproduct.proto\x12\aproduct\"[\n" +
	"\x15SearchProductsRequest\x12\x14\n" +
	"\x05query\x18\x01 \x01(\tR\x05query\x12\x14\n" +
	"\x05limit\x18\x02 \x01(\x05R\x05limit\x12\x16\n" +
	"\x06offset\x18\x03 \x01(\x05R\x06offset\"g\n" +
	"\x16SearchProductsResponse\x12,\n" +
	"\bproducts\x18\x01 \x03(\v2\x10.product.ProductR\bproducts\x12\x1f\n" +
	"\vtotal_count\x18\x02 \x01(\x05R\n" +
	"totalCount\"2\n" +
	"\x11GetProductRequest\x12\x1d\n" +
	"\n" +
	"product_id\x18\x01 \x01(\x05R\tproductId\"a\n" +
	"\x11AddProductRequest\x12\x14\n" +
	"\x05title\x18\x01 \x01(\tR\x05title\x12 \n" +
	"\vdescription\x18\x02 \x01(\tR\vdescription\x12\x14\n" +
	"\x05price\x18\x03 \x01(\x01R\x05price\"3\n" +
	"\x12AddProductResponse\x12\x1d\n" +
	"\n" +
	"product_id\x18\x01 \x01(\x05R\tproductId\"5\n" +
	"\x14DeleteProductRequest\x12\x1d\n" +
	"\n" +
	"product_id\x18\x01 \x01(\x05R\tproductId\"1\n" +
	"\x15DeleteProductResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\"g\n" +
	"\aProduct\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x05R\x02id\x12\x14\n" +
	"\x05title\x18\x02 \x01(\tR\x05title\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescription\x12\x14\n" +
	"\x05price\x18\x04 \x01(\x01R\x05price2\xb6\x02\n" +
	"\x0eProductService\x12:\n" +
	"\n" +
	"GetProduct\x12\x1a.product.GetProductRequest\x1a\x10.product.Product\x12Q\n" +
	"\x0eSearchProducts\x12\x1e.product.SearchProductsRequest\x1a\x1f.product.SearchProductsResponse\x12E\n" +
	"\n" +
	"AddProduct\x12\x1a.product.AddProductRequest\x1a\x1b.product.AddProductResponse\x12N\n" +
	"\rDeleteProduct\x12\x1d.product.DeleteProductRequest\x1a\x1e.product.DeleteProductResponseB\x0fZ\rproduct/protob\x06proto3"

var (
	file_product_proto_rawDescOnce sync.Once
	file_product_proto_rawDescData []byte
)

func file_product_proto_rawDescGZIP() []byte {
	file_product_proto_rawDescOnce.Do(func() {
		file_product_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_product_proto_rawDesc), len(file_product_proto_rawDesc)))
	})
	return file_product_proto_rawDescData
}

var file_product_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_product_proto_goTypes = []any{
	(*SearchProductsRequest)(nil),  // 0: product.SearchProductsRequest
	(*SearchProductsResponse)(nil), // 1: product.SearchProductsResponse
	(*GetProductRequest)(nil),      // 2: product.GetProductRequest
	(*AddProductRequest)(nil),      // 3: product.AddProductRequest
	(*AddProductResponse)(nil),     // 4: product.AddProductResponse
	(*DeleteProductRequest)(nil),   // 5: product.DeleteProductRequest
	(*DeleteProductResponse)(nil),  // 6: product.DeleteProductResponse
	(*Product)(nil),                // 7: product.Product
}
var file_product_proto_depIdxs = []int32{
	7, // 0: product.SearchProductsResponse.products:type_name -> product.Product
	2, // 1: product.ProductService.GetProduct:input_type -> product.GetProductRequest
	0, // 2: product.ProductService.SearchProducts:input_type -> product.SearchProductsRequest
	3, // 3: product.ProductService.AddProduct:input_type -> product.AddProductRequest
	5, // 4: product.ProductService.DeleteProduct:input_type -> product.DeleteProductRequest
	7, // 5: product.ProductService.GetProduct:output_type -> product.Product
	1, // 6: product.ProductService.SearchProducts:output_type -> product.SearchProductsResponse
	4, // 7: product.ProductService.AddProduct:output_type -> product.AddProductResponse
	6, // 8: product.ProductService.DeleteProduct:output_type -> product.DeleteProductResponse
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_product_proto_init() }
func file_product_proto_init() {
	if File_product_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_product_proto_rawDesc), len(file_product_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_product_proto_goTypes,
		DependencyIndexes: file_product_proto_depIdxs,
		MessageInfos:      file_product_proto_msgTypes,
	}.Build()
	File_product_proto = out.File
	file_product_proto_goTypes = nil
	file_product_proto_depIdxs = nil
}
