package svc

import (
	"webtranshome/internal/config"
	"webtranshome/internal/middleware"
	"webtranshome/internal/model"
	catalogcommand "webtranshome/internal/modules/catalog/application/command"
	catalogquery "webtranshome/internal/modules/catalog/application/query"
	catalogpostgres "webtranshome/internal/modules/catalog/infrastructure/postgres"
	"webtranshome/internal/modules/identityaccess/application/command"
	"webtranshome/internal/modules/identityaccess/application/passwordreset"
	"webtranshome/internal/modules/identityaccess/application/query"
	"webtranshome/internal/modules/identityaccess/infrastructure/notification"
	passwordbcrypt "webtranshome/internal/modules/identityaccess/infrastructure/password/bcrypt"
	authpostgres "webtranshome/internal/modules/identityaccess/infrastructure/postgres"
	jwtissuer "webtranshome/internal/modules/identityaccess/infrastructure/token/jwtissuer"
	mediaassets "webtranshome/internal/modules/media/application/assets"
	mediaassetspg "webtranshome/internal/modules/media/infrastructure/postgres"
	mediastorage "webtranshome/internal/modules/media/infrastructure/storage"

	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config                   config.Config
	UserModel                model.UsersModel
	LoginAdmin               *command.Login
	PasswordReset            *passwordreset.Service
	AccessControlQueries     *query.AccessControl
	ReplaceAccessAssignments *command.ReplaceAccessAssignments
	ManageRoles              *command.ManageRoles
	ListCategories           *catalogquery.ListCategories
	GetCategory              *catalogquery.GetCategory
	CreateCategory           *catalogcommand.CreateCategory
	UpdateCategory           *catalogcommand.UpdateCategory
	ListBrands               *catalogquery.ListBrands
	GetBrand                 *catalogquery.GetBrand
	CreateBrand              *catalogcommand.CreateBrand
	UpdateBrand              *catalogcommand.UpdateBrand
	ListRooms                *catalogquery.ListRooms
	GetRoom                  *catalogquery.GetRoom
	CreateRoom               *catalogcommand.CreateRoom
	UpdateRoom               *catalogcommand.UpdateRoom
	ListProducts             *catalogquery.ListProducts
	GetProduct               *catalogquery.GetProduct
	CreateProduct            *catalogcommand.CreateProduct
	UpdateProduct            *catalogcommand.UpdateProduct
	ReplaceProductRooms      *catalogcommand.ReplaceProductRooms
	ImageAssets              *mediaassets.Service
	Authenticate             rest.Middleware
	RequirePermission        func(string) rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Database.DataSource)
	userModel := model.NewUsersModel(conn)
	userRepository := authpostgres.NewUserRepository(conn, userModel)
	permissionRepository := authpostgres.NewPermissionRepository(conn)
	accessControlRepository := authpostgres.NewAccessControlRepository(conn)
	categoryRepository := catalogpostgres.NewCategoryRepository(conn)
	loginAdmin := command.NewLogin(
		userRepository,
		passwordbcrypt.NewVerifier(),
		jwtissuer.NewIssuer(c.Auth.AccessSecret, c.Auth.AccessExpire),
	)
	checkPermission := query.NewCheckPermission(permissionRepository)
	resetRepository := authpostgres.NewPasswordResetRepository(conn)
	whatsappSender := configuredSender(c.PasswordReset.WhatsAppGatewayURL, c.PasswordReset.WhatsAppGatewayToken, "whatsapp", c.Mode == "dev")
	emailSender := configuredSender(c.PasswordReset.EmailGatewayURL, c.PasswordReset.EmailGatewayToken, "email", c.Mode == "dev")
	resetService := passwordreset.New(resetRepository, c.PasswordReset.HMACSecret, whatsappSender, emailSender)
	tokenVerifier := jwtissuer.NewVerifier(c.Auth.AccessSecret)
	accessControlQueries := query.NewAccessControl(
		accessControlRepository,
		accessControlRepository,
		accessControlRepository,
	)
	replaceAccessAssignments := command.NewReplaceAccessAssignments(accessControlRepository)
	manageRoles := command.NewManageRoles(accessControlRepository)
	listCategories := catalogquery.NewListCategories(categoryRepository)
	getCategory := catalogquery.NewGetCategory(categoryRepository)
	createCategory := catalogcommand.NewCreateCategory(categoryRepository, categoryRepository)
	updateCategory := catalogcommand.NewUpdateCategory(categoryRepository, categoryRepository)
	brandRepository := catalogpostgres.NewBrandRepository(conn)
	listBrands := catalogquery.NewListBrands(brandRepository)
	getBrand := catalogquery.NewGetBrand(brandRepository)
	createBrand := catalogcommand.NewCreateBrand(brandRepository)
	updateBrand := catalogcommand.NewUpdateBrand(brandRepository, brandRepository)
	roomRepository := catalogpostgres.NewRoomRepository(conn)
	listRooms := catalogquery.NewListRooms(roomRepository)
	getRoom := catalogquery.NewGetRoom(roomRepository)
	createRoom := catalogcommand.NewCreateRoom(roomRepository)
	updateRoom := catalogcommand.NewUpdateRoom(roomRepository, roomRepository)
	productRepository := catalogpostgres.NewProductRepository(conn)
	listProducts := catalogquery.NewListProducts(productRepository)
	getProduct := catalogquery.NewGetProduct(productRepository)
	createProduct := catalogcommand.NewCreateProduct(productRepository, productRepository)
	updateProduct := catalogcommand.NewUpdateProduct(productRepository, productRepository)
	replaceProductRooms := catalogcommand.NewReplaceProductRooms(productRepository, roomRepository, productRepository)
	imageAssetRepository := mediaassetspg.NewAssetRepository(conn)
	imageAssets := mediaassets.New(imageAssetRepository, mediastorage.NewLocal(c.ImageStorage.Directory, c.ImageStorage.MaxBytes))

	return &ServiceContext{
		Config:                   c,
		UserModel:                userModel,
		LoginAdmin:               loginAdmin,
		PasswordReset:            resetService,
		AccessControlQueries:     accessControlQueries,
		ReplaceAccessAssignments: replaceAccessAssignments,
		ManageRoles:              manageRoles,
		ListCategories:           listCategories,
		GetCategory:              getCategory,
		CreateCategory:           createCategory,
		UpdateCategory:           updateCategory,
		ListBrands:               listBrands,
		GetBrand:                 getBrand,
		CreateBrand:              createBrand,
		UpdateBrand:              updateBrand,
		ListRooms:                listRooms,
		GetRoom:                  getRoom,
		CreateRoom:               createRoom,
		UpdateRoom:               updateRoom,
		ListProducts:             listProducts,
		GetProduct:               getProduct,
		CreateProduct:            createProduct,
		UpdateProduct:            updateProduct,
		ReplaceProductRooms:      replaceProductRooms,
		ImageAssets:              imageAssets,
		Authenticate:             middleware.Authenticate(tokenVerifier, userRepository),
		RequirePermission: func(code string) rest.Middleware {
			return middleware.RequirePermission(checkPermission, code)
		},
	}
}

func configuredSender(endpoint, token, channel string, development bool) notification.Sender {
	if endpoint != "" {
		return notification.NewGatewaySender(endpoint, token)
	}
	if development {
		return notification.RedactedDevelopmentSender{Channel: channel}
	}
	return nil
}
