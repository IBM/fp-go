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

package function

// CurryAB converts a binary function f(A, B) into a curried function A → B → R.
//
// The suffix "AB" reflects the order in which the underlying function receives
// its arguments: first A, then B (i.e. the natural order is preserved).
//
// Type Parameters:
//   - FCT: The binary function type, constrained to func(A, B) R
//   - A: The type of the first parameter
//   - B: The type of the second parameter
//   - R: The return type
//
// Parameters:
//   - f: The binary function to curry
//
// Returns:
//   - A curried function that takes A and returns a function from B to R
//
// See Also:
//   - CurryBA: Curries a binary function with swapped argument order
func CurryAB[FCT ~func(A, B) R, A, B, R any](f FCT) func(A) func(B) R {
	return func(a A) func(B) R {
		return func(b B) R {
			return f(a, b)
		}
	}
}

// CurryBA converts a binary function f(B, A) into a curried function A → B → R.
//
// The suffix "BA" reflects the order in which the underlying function receives
// its arguments: first B, then A. The curried signature accepts A first and B
// second, so the arguments are swapped before forwarding to f.
//
// Type Parameters:
//   - FCT: The binary function type, constrained to func(B, A) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The binary function to curry (takes B before A)
//
// Returns:
//   - A curried function that takes A and returns a function from B to R,
//     forwarding the arguments to f as f(b, a)
//
// See Also:
//   - CurryAB: Curries a binary function preserving argument order
func CurryBA[FCT ~func(B, A) R, A, B, R any](f FCT) func(A) func(B) R {
	return func(a A) func(B) R {
		return func(b B) R {
			return f(b, a)
		}
	}
}

// CurryABC converts a ternary function f(A, B, C) into a curried function A → B → C → R.
//
// The suffix "ABC" reflects the order in which the underlying function receives
// its arguments: A, B, C (the natural order is preserved).
//
// Type Parameters:
//   - FCT: The ternary function type, constrained to func(A, B, C) R
//   - A: The type of the first parameter
//   - B: The type of the second parameter
//   - C: The type of the third parameter
//   - R: The return type
//
// Parameters:
//   - f: The ternary function to curry
//
// Returns:
//   - A curried function A → B → C → R
func CurryABC[FCT ~func(A, B, C) R, A, B, C, R any](f FCT) func(A) func(B) func(C) R {
	return func(a A) func(B) func(C) R {
		return func(b B) func(C) R {
			return func(c C) R {
				return f(a, b, c)
			}
		}
	}
}

// CurryACB converts a ternary function f(A, C, B) into a curried function A → B → C → R.
//
// The suffix "ACB" reflects the order in which the underlying function receives
// its arguments: A, C, B. The curried signature accepts them as A, B, C and
// forwards them to f as f(a, c, b).
//
// Type Parameters:
//   - FCT: The ternary function type, constrained to func(A, C, B) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The ternary function to curry (takes arguments in A, C, B order)
//
// Returns:
//   - A curried function A → B → C → R, forwarding to f as f(a, c, b)
func CurryACB[FCT ~func(A, C, B) R, A, B, C, R any](f FCT) func(A) func(B) func(C) R {
	return func(a A) func(B) func(C) R {
		return func(b B) func(C) R {
			return func(c C) R {
				return f(a, c, b)
			}
		}
	}
}

// CurryBAC converts a ternary function f(B, A, C) into a curried function A → B → C → R.
//
// The suffix "BAC" reflects the order in which the underlying function receives
// its arguments: B, A, C. The curried signature accepts them as A, B, C and
// forwards them to f as f(b, a, c).
//
// Type Parameters:
//   - FCT: The ternary function type, constrained to func(B, A, C) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The ternary function to curry (takes arguments in B, A, C order)
//
// Returns:
//   - A curried function A → B → C → R, forwarding to f as f(b, a, c)
func CurryBAC[FCT ~func(B, A, C) R, A, B, C, R any](f FCT) func(A) func(B) func(C) R {
	return func(a A) func(B) func(C) R {
		return func(b B) func(C) R {
			return func(c C) R {
				return f(b, a, c)
			}
		}
	}
}

// CurryBCA converts a ternary function f(B, C, A) into a curried function A → B → C → R.
//
// The suffix "BCA" reflects the order in which the underlying function receives
// its arguments: B, C, A. The curried signature accepts them as A, B, C and
// forwards them to f as f(b, c, a).
//
// Type Parameters:
//   - FCT: The ternary function type, constrained to func(B, C, A) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The ternary function to curry (takes arguments in B, C, A order)
//
// Returns:
//   - A curried function A → B → C → R, forwarding to f as f(b, c, a)
func CurryBCA[FCT ~func(B, C, A) R, A, B, C, R any](f FCT) func(A) func(B) func(C) R {
	return func(a A) func(B) func(C) R {
		return func(b B) func(C) R {
			return func(c C) R {
				return f(b, c, a)
			}
		}
	}
}

// CurryCAB converts a ternary function f(C, A, B) into a curried function A → B → C → R.
//
// The suffix "CAB" reflects the order in which the underlying function receives
// its arguments: C, A, B. The curried signature accepts them as A, B, C and
// forwards them to f as f(c, a, b).
//
// Type Parameters:
//   - FCT: The ternary function type, constrained to func(C, A, B) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The ternary function to curry (takes arguments in C, A, B order)
//
// Returns:
//   - A curried function A → B → C → R, forwarding to f as f(c, a, b)
func CurryCAB[FCT ~func(C, A, B) R, A, B, C, R any](f FCT) func(A) func(B) func(C) R {
	return func(a A) func(B) func(C) R {
		return func(b B) func(C) R {
			return func(c C) R {
				return f(c, a, b)
			}
		}
	}
}

// CurryCBA converts a ternary function f(C, B, A) into a curried function A → B → C → R.
//
// The suffix "CBA" reflects the order in which the underlying function receives
// its arguments: C, B, A (the reverse of natural order). The curried signature
// accepts them as A, B, C and forwards them to f as f(c, b, a).
//
// Type Parameters:
//   - FCT: The ternary function type, constrained to func(C, B, A) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The ternary function to curry (takes arguments in C, B, A order)
//
// Returns:
//   - A curried function A → B → C → R, forwarding to f as f(c, b, a)
func CurryCBA[FCT ~func(C, B, A) R, A, B, C, R any](f FCT) func(A) func(B) func(C) R {
	return func(a A) func(B) func(C) R {
		return func(b B) func(C) R {
			return func(c C) R {
				return f(c, b, a)
			}
		}
	}
}

// CurryABCD converts a quaternary function f(A, B, C, D) into a curried function A → B → C → D → R.
//
// The suffix "ABCD" reflects the order in which the underlying function receives
// its arguments: A, B, C, D (the natural order is preserved).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(A, B, C, D) R
//   - A: The type of the first parameter
//   - B: The type of the second parameter
//   - C: The type of the third parameter
//   - D: The type of the fourth parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(a, b, c, d)
func CurryABCD[FCT ~func(A, B, C, D) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(a, b, c, d)
				}
			}
		}
	}
}

// CurryABDC converts a quaternary function f(A, B, D, C) into a curried function A → B → C → D → R.
//
// The suffix "ABDC" reflects the order in which the underlying function receives
// its arguments: A, B, D, C. The curried signature accepts them as A, B, C, D
// and forwards to f as f(a, b, d, c).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(A, B, D, C) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in A, B, D, C order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(a, b, d, c)
func CurryABDC[FCT ~func(A, B, D, C) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(a, b, d, c)
				}
			}
		}
	}
}

// CurryACBD converts a quaternary function f(A, C, B, D) into a curried function A → B → C → D → R.
//
// The suffix "ACBD" reflects the order in which the underlying function receives
// its arguments: A, C, B, D. The curried signature accepts them as A, B, C, D
// and forwards to f as f(a, c, b, d).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(A, C, B, D) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in A, C, B, D order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(a, c, b, d)
func CurryACBD[FCT ~func(A, C, B, D) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(a, c, b, d)
				}
			}
		}
	}
}

// CurryACDB converts a quaternary function f(A, C, D, B) into a curried function A → B → C → D → R.
//
// The suffix "ACDB" reflects the order in which the underlying function receives
// its arguments: A, C, D, B. The curried signature accepts them as A, B, C, D
// and forwards to f as f(a, c, d, b).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(A, C, D, B) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in A, C, D, B order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(a, c, d, b)
func CurryACDB[FCT ~func(A, C, D, B) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(a, c, d, b)
				}
			}
		}
	}
}

// CurryADBC converts a quaternary function f(A, D, B, C) into a curried function A → B → C → D → R.
//
// The suffix "ADBC" reflects the order in which the underlying function receives
// its arguments: A, D, B, C. The curried signature accepts them as A, B, C, D
// and forwards to f as f(a, d, b, c).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(A, D, B, C) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in A, D, B, C order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(a, d, b, c)
func CurryADBC[FCT ~func(A, D, B, C) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(a, d, b, c)
				}
			}
		}
	}
}

// CurryADCB converts a quaternary function f(A, D, C, B) into a curried function A → B → C → D → R.
//
// The suffix "ADCB" reflects the order in which the underlying function receives
// its arguments: A, D, C, B. The curried signature accepts them as A, B, C, D
// and forwards to f as f(a, d, c, b).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(A, D, C, B) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in A, D, C, B order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(a, d, c, b)
func CurryADCB[FCT ~func(A, D, C, B) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(a, d, c, b)
				}
			}
		}
	}
}

// CurryBACD converts a quaternary function f(B, A, C, D) into a curried function A → B → C → D → R.
//
// The suffix "BACD" reflects the order in which the underlying function receives
// its arguments: B, A, C, D. The curried signature accepts them as A, B, C, D
// and forwards to f as f(b, a, c, d).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(B, A, C, D) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in B, A, C, D order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(b, a, c, d)
func CurryBACD[FCT ~func(B, A, C, D) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(b, a, c, d)
				}
			}
		}
	}
}

// CurryBADC converts a quaternary function f(B, A, D, C) into a curried function A → B → C → D → R.
//
// The suffix "BADC" reflects the order in which the underlying function receives
// its arguments: B, A, D, C. The curried signature accepts them as A, B, C, D
// and forwards to f as f(b, a, d, c).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(B, A, D, C) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in B, A, D, C order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(b, a, d, c)
func CurryBADC[FCT ~func(B, A, D, C) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(b, a, d, c)
				}
			}
		}
	}
}

// CurryBCAD converts a quaternary function f(B, C, A, D) into a curried function A → B → C → D → R.
//
// The suffix "BCAD" reflects the order in which the underlying function receives
// its arguments: B, C, A, D. The curried signature accepts them as A, B, C, D
// and forwards to f as f(b, c, a, d).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(B, C, A, D) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in B, C, A, D order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(b, c, a, d)
func CurryBCAD[FCT ~func(B, C, A, D) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(b, c, a, d)
				}
			}
		}
	}
}

// CurryBCDA converts a quaternary function f(B, C, D, A) into a curried function A → B → C → D → R.
//
// The suffix "BCDA" reflects the order in which the underlying function receives
// its arguments: B, C, D, A. The curried signature accepts them as A, B, C, D
// and forwards to f as f(b, c, d, a).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(B, C, D, A) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in B, C, D, A order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(b, c, d, a)
func CurryBCDA[FCT ~func(B, C, D, A) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(b, c, d, a)
				}
			}
		}
	}
}

// CurryBDAC converts a quaternary function f(B, D, A, C) into a curried function A → B → C → D → R.
//
// The suffix "BDAC" reflects the order in which the underlying function receives
// its arguments: B, D, A, C. The curried signature accepts them as A, B, C, D
// and forwards to f as f(b, d, a, c).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(B, D, A, C) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in B, D, A, C order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(b, d, a, c)
func CurryBDAC[FCT ~func(B, D, A, C) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(b, d, a, c)
				}
			}
		}
	}
}

// CurryBDCA converts a quaternary function f(B, D, C, A) into a curried function A → B → C → D → R.
//
// The suffix "BDCA" reflects the order in which the underlying function receives
// its arguments: B, D, C, A. The curried signature accepts them as A, B, C, D
// and forwards to f as f(b, d, c, a).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(B, D, C, A) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in B, D, C, A order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(b, d, c, a)
func CurryBDCA[FCT ~func(B, D, C, A) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(b, d, c, a)
				}
			}
		}
	}
}

// CurryCABD converts a quaternary function f(C, A, B, D) into a curried function A → B → C → D → R.
//
// The suffix "CABD" reflects the order in which the underlying function receives
// its arguments: C, A, B, D. The curried signature accepts them as A, B, C, D
// and forwards to f as f(c, a, b, d).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(C, A, B, D) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in C, A, B, D order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(c, a, b, d)
func CurryCABD[FCT ~func(C, A, B, D) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(c, a, b, d)
				}
			}
		}
	}
}

// CurryCADB converts a quaternary function f(C, A, D, B) into a curried function A → B → C → D → R.
//
// The suffix "CADB" reflects the order in which the underlying function receives
// its arguments: C, A, D, B. The curried signature accepts them as A, B, C, D
// and forwards to f as f(c, a, d, b).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(C, A, D, B) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in C, A, D, B order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(c, a, d, b)
func CurryCADB[FCT ~func(C, A, D, B) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(c, a, d, b)
				}
			}
		}
	}
}

// CurryCBAD converts a quaternary function f(C, B, A, D) into a curried function A → B → C → D → R.
//
// The suffix "CBAD" reflects the order in which the underlying function receives
// its arguments: C, B, A, D. The curried signature accepts them as A, B, C, D
// and forwards to f as f(c, b, a, d).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(C, B, A, D) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in C, B, A, D order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(c, b, a, d)
func CurryCBAD[FCT ~func(C, B, A, D) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(c, b, a, d)
				}
			}
		}
	}
}

// CurryCBDA converts a quaternary function f(C, B, D, A) into a curried function A → B → C → D → R.
//
// The suffix "CBDA" reflects the order in which the underlying function receives
// its arguments: C, B, D, A. The curried signature accepts them as A, B, C, D
// and forwards to f as f(c, b, d, a).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(C, B, D, A) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in C, B, D, A order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(c, b, d, a)
func CurryCBDA[FCT ~func(C, B, D, A) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(c, b, d, a)
				}
			}
		}
	}
}

// CurryCDAB converts a quaternary function f(C, D, A, B) into a curried function A → B → C → D → R.
//
// The suffix "CDAB" reflects the order in which the underlying function receives
// its arguments: C, D, A, B. The curried signature accepts them as A, B, C, D
// and forwards to f as f(c, d, a, b).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(C, D, A, B) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in C, D, A, B order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(c, d, a, b)
func CurryCDAB[FCT ~func(C, D, A, B) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(c, d, a, b)
				}
			}
		}
	}
}

// CurryCDBA converts a quaternary function f(C, D, B, A) into a curried function A → B → C → D → R.
//
// The suffix "CDBA" reflects the order in which the underlying function receives
// its arguments: C, D, B, A. The curried signature accepts them as A, B, C, D
// and forwards to f as f(c, d, b, a).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(C, D, B, A) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in C, D, B, A order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(c, d, b, a)
func CurryCDBA[FCT ~func(C, D, B, A) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(c, d, b, a)
				}
			}
		}
	}
}

// CurryDABC converts a quaternary function f(D, A, B, C) into a curried function A → B → C → D → R.
//
// The suffix "DABC" reflects the order in which the underlying function receives
// its arguments: D, A, B, C. The curried signature accepts them as A, B, C, D
// and forwards to f as f(d, a, b, c).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(D, A, B, C) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in D, A, B, C order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(d, a, b, c)
func CurryDABC[FCT ~func(D, A, B, C) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(d, a, b, c)
				}
			}
		}
	}
}

// CurryDACB converts a quaternary function f(D, A, C, B) into a curried function A → B → C → D → R.
//
// The suffix "DACB" reflects the order in which the underlying function receives
// its arguments: D, A, C, B. The curried signature accepts them as A, B, C, D
// and forwards to f as f(d, a, c, b).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(D, A, C, B) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in D, A, C, B order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(d, a, c, b)
func CurryDACB[FCT ~func(D, A, C, B) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(d, a, c, b)
				}
			}
		}
	}
}

// CurryDBAC converts a quaternary function f(D, B, A, C) into a curried function A → B → C → D → R.
//
// The suffix "DBAC" reflects the order in which the underlying function receives
// its arguments: D, B, A, C. The curried signature accepts them as A, B, C, D
// and forwards to f as f(d, b, a, c).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(D, B, A, C) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in D, B, A, C order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(d, b, a, c)
func CurryDBAC[FCT ~func(D, B, A, C) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(d, b, a, c)
				}
			}
		}
	}
}

// CurryDBCA converts a quaternary function f(D, B, C, A) into a curried function A → B → C → D → R.
//
// The suffix "DBCA" reflects the order in which the underlying function receives
// its arguments: D, B, C, A. The curried signature accepts them as A, B, C, D
// and forwards to f as f(d, b, c, a).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(D, B, C, A) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in D, B, C, A order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(d, b, c, a)
func CurryDBCA[FCT ~func(D, B, C, A) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(d, b, c, a)
				}
			}
		}
	}
}

// CurryDCAB converts a quaternary function f(D, C, A, B) into a curried function A → B → C → D → R.
//
// The suffix "DCAB" reflects the order in which the underlying function receives
// its arguments: D, C, A, B. The curried signature accepts them as A, B, C, D
// and forwards to f as f(d, c, a, b).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(D, C, A, B) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in D, C, A, B order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(d, c, a, b)
func CurryDCAB[FCT ~func(D, C, A, B) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(d, c, a, b)
				}
			}
		}
	}
}

// CurryDCBA converts a quaternary function f(D, C, B, A) into a curried function A → B → C → D → R.
//
// The suffix "DCBA" reflects the order in which the underlying function receives
// its arguments: D, C, B, A (the reverse of natural order). The curried signature
// accepts them as A, B, C, D and forwards to f as f(d, c, b, a).
//
// Type Parameters:
//   - FCT: The quaternary function type, constrained to func(D, C, B, A) R
//   - A: The type of the first curried parameter
//   - B: The type of the second curried parameter
//   - C: The type of the third curried parameter
//   - D: The type of the fourth curried parameter
//   - R: The return type
//
// Parameters:
//   - f: The quaternary function to curry (takes arguments in D, C, B, A order)
//
// Returns:
//   - A curried function A → B → C → D → R, forwarding to f as f(d, c, b, a)
func CurryDCBA[FCT ~func(D, C, B, A) R, A, B, C, D, R any](f FCT) func(A) func(B) func(C) func(D) R {
	return func(a A) func(B) func(C) func(D) R {
		return func(b B) func(C) func(D) R {
			return func(c C) func(D) R {
				return func(d D) R {
					return f(d, c, b, a)
				}
			}
		}
	}
}
