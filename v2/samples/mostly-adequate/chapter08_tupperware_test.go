// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mostlyadequate

import (
	"fmt"
	"time"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	I "github.com/IBM/fp-go/v2/identity"
	"github.com/IBM/fp-go/v2/ioresult"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/ord"
	R "github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
)

type Account struct {
	Balance float32
}

func MakeAccount(b float32) Account {
	return Account{Balance: b}
}

func getBalance(a Account) float32 {
	return a.Balance
}

type (
	Chapter08User struct {
		Id     int
		Name   string
		Active bool
		Saved  bool
	}
)

var (
	albert08 = Chapter08User{
		Id:     1,
		Active: true,
		Name:   "Albert",
	}

	gary08 = Chapter08User{
		Id:     2,
		Active: false,
		Name:   "Gary",
	}

	theresa08 = Chapter08User{
		Id:     3,
		Active: true,
		Name:   "Theresa",
	}

	yi08 = Chapter08User{Id: 4, Name: "Yi", Active: true}
)

func (u Chapter08User) getName() string {
	return u.Name
}

func (u Chapter08User) isActive() bool {
	return u.Active
}

var (
	ordFloat32       = ord.FromStrictCompare[float32]()
	UpdateLedger     = F.Identity[Account]
	RemainingBalance = F.Flow2(
		getBalance,
		S.Format[float32]("Your balance is $%0.2f"),
	)
	FinishTransaction = F.Flow2(
		UpdateLedger,
		RemainingBalance,
	)
	getTwenty = F.Flow2(
		Withdraw(20),
		O.Fold(F.Constant("You're broke!"), FinishTransaction),
	)

	// showWelcome :: User -> String
	showWelcome = F.Flow2(
		Chapter08User.getName,
		S.Format[string]("Welcome %s"),
	)

	// checkActive :: User -> Either error User
	checkActive = R.FromPredicate(Chapter08User.isActive, F.Constant1[Chapter08User](fmt.Errorf("your account is not active")))

	// validateUser :: (User -> Either String ()) -> User -> Either String User
	validateUser = F.Curry2(func(validate func(Chapter08User) Result[any], user Chapter08User) Result[Chapter08User] {
		return F.Pipe2(
			user,
			validate,
			R.MapTo[any](user),
		)
	})

	// save :: User -> IOEither error User
	save = func(user Chapter08User) IOResult[Chapter08User] {
		return ioresult.FromIO(func() Chapter08User {
			var u = user
			u.Saved = true
			return u
		})
	}
)

// Withdraw creates a Kleisli arrow that attempts to withdraw an amount from an account.
// Returns Some(account) if sufficient balance, None otherwise.
// This demonstrates the Option Kleisli type in action.
func Withdraw(amount float32) O.Kleisli[Account, Account] {
	return F.Flow3(
		getBalance,
		O.FromPredicate(ord.Geq(ordFloat32)(amount)),
		O.Map(F.Flow2(
			N.Add(-amount),
			MakeAccount,
		)))
}

type User struct {
	BirthDate string
}

func getBirthDate(u User) string {
	return u.BirthDate
}

func MakeUser(d string) User {
	return User{BirthDate: d}
}

var parseDate = F.Bind1of2(R.Eitherize2(time.Parse))(time.DateOnly)

// GetAge creates a Result Kleisli arrow that calculates age in days from a User's birth date.
// This demonstrates the Result Kleisli type for computations that may fail.
func GetAge(now time.Time) R.Kleisli[User, float64] {
	return F.Flow3(
		getBirthDate,
		parseDate,
		R.Map(F.Flow3(
			now.Sub,
			time.Duration.Hours,
			N.Mul(1/24.0),
		)),
	)
}

func Example_widthdraw() {
	fmt.Println(getTwenty(MakeAccount(200)))
	fmt.Println(getTwenty(MakeAccount(10)))

	// Output:
	// Your balance is $180.00
	// You're broke!
}

func Example_getAge() {
	now, err := time.Parse(time.DateOnly, "2023-09-01")
	if err != nil {
		panic(err)
	}

	fmt.Println(GetAge(now)(MakeUser("2005-12-12")))
	fmt.Println(GetAge(now)(MakeUser("July 4, 2001")))

	fortune := F.Flow3(
		N.Add(365.0),
		S.Format[float64]("%0.0f"),
		Concat("If you survive, you will be "),
	)

	zoltar := F.Flow3(
		GetAge(now),
		R.Map(fortune),
		R.GetOrElse(errors.ToString),
	)

	fmt.Println(zoltar(MakeUser("2005-12-12")))

	// Output:
	// Right[float64](6472)
	// Left[*time.ParseError](parsing time "July 4, 2001" as "2006-01-02": cannot parse "July 4, 2001" as "2006")
	// If you survive, you will be 6837
}

func Example_solution08A() {
	incrF := I.Map(N.Add(1))

	fmt.Println(incrF(I.Of(2)))

	// Output: 3
}

func Example_solution08B() {
	// initial :: User -> Option rune
	initial := F.Flow3(
		Chapter08User.getName,
		S.ToRunes,
		A.Head[rune],
	)

	fmt.Println(initial(albert08))

	// Output:
	// Some[int32](65)
}

func Example_solution08C() {

	// eitherWelcome :: User -> Either String String
	eitherWelcome := F.Flow2(
		checkActive,
		R.Map(showWelcome),
	)

	fmt.Println(eitherWelcome(gary08))
	fmt.Println(eitherWelcome(theresa08))

	// Output:
	// Left[*errors.errorString](your account is not active)
	// Right[string](Welcome Theresa)
}

func Example_solution08D() {

	// // validateName :: User -> Either String ()
	validateName := F.Flow3(
		Chapter08User.getName,
		R.FromPredicate(F.Flow2(
			S.Size,
			ord.Gt(ord.FromStrictCompare[int]())(3),
		), errors.OnSome[string]("Your name %s is larger than 3 characters")),
		R.Map(F.ToAny[string]),
	)

	saveAndWelcome := F.Flow2(
		save,
		ioresult.Map(showWelcome),
	)

	register := F.Flow3(
		validateUser(validateName),
		ioresult.FromEither[Chapter08User],
		ioresult.Chain(saveAndWelcome),
	)

	fmt.Println(validateName(gary08))
	fmt.Println(validateName(yi08))

	fmt.Println(register(albert08)())
	fmt.Println(register(yi08)())

	// Output:
	// Right[string](Gary)
	// Left[*errors.errorString](Your name Yi is larger than 3 characters)
	// Right[string](Welcome Albert)
	// Left[*errors.errorString](Your name Yi is larger than 3 characters)
}
