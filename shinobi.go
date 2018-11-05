package main

import (
    "os"
    "fmt"
    "flag"
    "encoding/binary"
    "crypto/rand"
    "strconv"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
    "github.com/aws/aws-sdk-go/aws/credentials/stscreds"
    "github.com/aws/aws-sdk-go/service/sts"

    "golang.org/x/crypto/ssh/terminal"
    "github.com/olekukonko/tablewriter"
    "github.com/fatih/color"

    "github.com/joho/godotenv"
)

const (
    AppVersion = "0.0.4"
)

var (
    argProfile = flag.String("profile", "", "Profile 名を指定.")
    argRole = flag.String("role", "", "Role ARN を指定.")
    argRegion = flag.String("region", "ap-northeast-1", "Region 名を指定.")
    argEndpoint = flag.String("endpoint", "", "AWS API のエンドポイントを指定.")
    argVersion = flag.Bool("version", false, "バージョンを出力.")
    argCreate = flag.Bool("create", false, "ユーザーを作成.")
    argDelete = flag.Bool("delete", false, "ユーザーを削除.")
    argList = flag.Bool("list", false, "ユーザー一覧を取得.")
    argUsername = flag.String("username", "", "User Pool に作成するユーザー名を指定.")
    argPassword = flag.String("password", "", "User Pool に作成するユーザーのパスワードを指定.")
    argEmail = flag.String("email", "", "User Pool に作成するユーザーのメールアドレスを指定.")
    argNickname = flag.String("nickname", "", "User Pool に作成するユーザーのニックネームを指定.")

    poolId = ""
    clientId = ""
    cognitoClient *cognitoidentityprovider.CognitoIdentityProvider
)

func awsCognitoClient(profile string, region string, role string) *cognitoidentityprovider.CognitoIdentityProvider {
    var config aws.Config
    if profile != "" && role == "" {
        creds := credentials.NewSharedCredentials("", profile)
        config = aws.Config{Region: aws.String(region),
                            Credentials: creds,
                            Endpoint: aws.String(*argEndpoint)}
    } else if profile == "" && role != "" {
        sess := session.Must(session.NewSession())
        creds := stscreds.NewCredentials(sess, role)
        config = aws.Config{Region: aws.String(region),
            Credentials: creds,
            Endpoint: aws.String(*argEndpoint)}
    } else if profile != "" && role != "" {
        sess := session.Must(session.NewSessionWithOptions(session.Options{Profile:profile}))
        assumeRoler := sts.New(sess)
        creds := stscreds.NewCredentialsWithClient(assumeRoler, role)
        config = aws.Config{Region: aws.String(region),
            Credentials: creds,
            Endpoint: aws.String(*argEndpoint)}
    } else {
        config = aws.Config{Region: aws.String(region),
                            Endpoint: aws.String(*argEndpoint)}
    }

    sess := session.New(&config)
    cognitoClient := cognitoidentityprovider.New(sess)
    return cognitoClient
}

func getSession(poolId string, clientId string, userName string, tempPassword string) *string {
    params := &cognitoidentityprovider.AdminInitiateAuthInput {
        UserPoolId: aws.String(poolId),
        ClientId: aws.String(clientId),
        AuthFlow: aws.String("ADMIN_NO_SRP_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(userName),
			"PASSWORD": aws.String(tempPassword),
		},
    }
    
    res, err := cognitoClient.AdminInitiateAuth(params)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    session := res.Session
    return session
}

func generateTemporaryPassword() string {
    var n uint64
    binary.Read(rand.Reader, binary.LittleEndian, &n)
    tempPassword := strconv.FormatUint(n, 36)
    return tempPassword
}

func createUser(userName string, userPassword string, userEmail string, userNickname string) {
    var tempPassword string
    tempPassword = generateTemporaryPassword()
    params1 := &cognitoidentityprovider.AdminCreateUserInput {
        TemporaryPassword: aws.String(tempPassword),
        UserPoolId: aws.String(poolId),
        Username: aws.String(userName),
        UserAttributes: []*cognitoidentityprovider.AttributeType{
            {
                Name: aws.String("email"),
                Value: aws.String(userEmail),
            },
            {
                Name: aws.String("nickname"),
                Value: aws.String(userNickname),
            },
        },
    }

    _, err := cognitoClient.AdminCreateUser(params1)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    session := getSession(poolId, clientId, userName, tempPassword)

	params2 := &cognitoidentityprovider.AdminRespondToAuthChallengeInput{
	    ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
	    ChallengeResponses: map[string]*string{
	        "NEW_PASSWORD": aws.String(userPassword),
	        "USERNAME": aws.String(userName),
	    },
        ClientId: aws.String(clientId),
        Session: session,
        UserPoolId: aws.String(poolId),
	}

    _, err = cognitoClient.AdminRespondToAuthChallenge(params2)
	if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    } else {
        fmt.Println("ユーザー " + userName + " を作成しました.")
        listUsers(userName)
        os.Exit(0)
    }
}

func deleteUser(userName string) {
    params := &cognitoidentityprovider.AdminDeleteUserInput {
        UserPoolId: aws.String(poolId),
        Username: aws.String(userName),
    }

    _, err := cognitoClient.AdminDeleteUser(params)
	if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    } else {
        fmt.Println("ユーザー " + userName + " を削除しました.")
        os.Exit(0)
    }
}

func convertDate(d time.Time) (convertedDate string) {
    const layout = "2006-01-02 15:04:05"
    jst := time.FixedZone("Asia/Tokyo", 9*60*60)
    convertedDate = d.In(jst).Format(layout)

    return convertedDate
}

func outputTbl(data [][]string) {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Username", "Nickname", "Email", "UserStatus", "UserCreateDate", "UserLastModifiedDate"})
    for _, value := range data {
        table.Append(value)
    }
    table.Render()
}

func listUsers(userName string) {
    params := &cognitoidentityprovider.ListUsersInput {
        UserPoolId: aws.String(poolId),
    }

    if userName != "" {
        params.SetFilter("username = \"" + userName + "\"")
    }

    res, err := cognitoClient.ListUsers(params)
	if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    allUsers := [][]string{}
    for _, r := range res.Users {
        var nickname string
        var email string
        for _, a := range r.Attributes {
            if *a.Name == "nickname" {
                nickname = *a.Value
            } else if *a.Name == "email" {
                email = *a.Value
            }
        }
        createDate := convertDate(*r.UserCreateDate)
        modifiedDate := convertDate(*r.UserLastModifiedDate)
        User := []string{
            *r.Username,
            nickname,
            email,
            *r.UserStatus,
            createDate,
            modifiedDate,
        }
        allUsers = append(allUsers, User)
    }

    userPoolName := getUserPoolName()
    fmt.Printf("%v %v %v %v\n", color.HiGreenString("User Pool Name:"), userPoolName,
                                color.HiGreenString("User Pool ID:"), poolId)
    outputTbl(allUsers)
}

func getUserPoolName() string {
    params := &cognitoidentityprovider.DescribeUserPoolInput {
        UserPoolId: aws.String(poolId),
    }

    res, err := cognitoClient.DescribeUserPool(params)
	if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    var userPoolName string
    userPoolName = *res.UserPool.Name
    return userPoolName
}

func main() {
    flag.Parse()

    if *argVersion {
        fmt.Println(AppVersion)
        os.Exit(0)
    }

    // .env ファイルから COGNITO_USER_POOL_ID と COGNITO_CLIENT_ID を読み取る.
    err := godotenv.Load()
    if err != nil {
        fmt.Println(".env ファイルの読み込みに失敗しました.")
    }

    if os.Getenv("COGNITO_USER_POOL_ID") == "" {
        fmt.Println("Cognito User Pool ID (環境変数: COGNITO_USER_POOL_ID) を指定して下さい.")
        os.Exit(1)
    } else {
        poolId = os.Getenv("COGNITO_USER_POOL_ID")
    }

    if os.Getenv("COGNITO_CLIENT_ID") == "" {
        fmt.Println("Cognito User Pool アプリクライアント ID (環境変数: COGNITO_CLIENT_ID) を指定して下さい.")
        os.Exit(1)
    } else {
        clientId = os.Getenv("COGNITO_CLIENT_ID")
    }

    cognitoClient = awsCognitoClient(*argProfile, *argRegion, *argRole)

    if *argCreate {
        if *argUsername == "" {
            fmt.Println("ユーザー名を指定して下さい.")
            os.Exit(1)
        }

        if *argEmail == "" {
            fmt.Println("メールアドレスを指定して下さい.")
            os.Exit(1)
        }
        if *argNickname == "" {
            fmt.Println("ニックネームを指定して下さい.")
            os.Exit(1)
        }

        var passwordValue string
        if *argPassword == "" {
            fmt.Println("パスワードを入力して下さい: ")
            maskedValue1, err := terminal.ReadPassword(0)
            if err != nil {
                fmt.Println("入力した値が不正です.")
                os.Exit(1)
            }
            fmt.Println("パスワードをもう一度入力して下さい: ")
            maskedValue2, err := terminal.ReadPassword(0)
            if err == nil && string(maskedValue1) == string(maskedValue2) {
                passwordValue = string(maskedValue2)
            } else {
                fmt.Println("入力した値が不正です.")
                os.Exit(1)
            }
        } else {
            passwordValue = *argPassword
        }
        createUser(*argUsername, passwordValue, *argEmail, *argNickname)
    } else if *argDelete {
        listUsers(*argUsername)
        fmt.Print("上記のユーザーを削除しますか?(y/n): ")
        var stdin string
        fmt.Scan(&stdin)
        switch stdin {
        case "y", "Y":
            fmt.Println("ユーザーを削除します.")
            deleteUser(*argUsername)
        case "n", "N":
            fmt.Println("処理を停止します.")
            os.Exit(0)
        default:
            fmt.Println("処理を停止します.")
            os.Exit(0)
        }
        deleteUser(*argUsername)
    } else {
        listUsers("")
    }
}
