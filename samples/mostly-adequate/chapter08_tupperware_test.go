// Copyright (c) 2023 IBM Corp.
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

	E "github.com/IBM/fp-go/either"
	"github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	N "github.com/IBM/fp-go/number"
	O "github.com/IBM/fp-go/option"
	"github.com/IBM/fp-go/ord"
	S "github.com/IBM/fp-go/string"
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
)

func Withdraw(amount float32) func(account Account) O.Option[Account] {

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

var parseDate = F.Bind1of2(E.Eitherize2(time.Parse))(time.DateOnly)

func GetAge(now time.Time) func(User) E.Either[error, float64] {
	return F.Flow3(
		getBirthDate,
		parseDate,
		E.Map[error](F.Flow3(
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
		E.Map[error](fortune),
		E.GetOrElse(errors.ToString),
	)

	fmt.Println(zoltar(MakeUser("2005-12-12")))

	// Output:
	// Right[<nil>, float64](6472)
	// Left[*time.ParseError, float64](parsing time "July 4, 2001" as "2006-01-02": cannot parse "July 4, 2001" as "2006")
	// If you survive, you will be 6837
}
