// Package codec provides a functional approach to encoding and decoding data with validation.
//
// The codec package combines the concepts of encoders and decoders into a unified Type that can
// both encode values to an output format and decode/validate values from an input format. This
// is particularly useful for data serialization, API validation, and type-safe transformations.
//
// # Core Concepts
//
// Type[A, O, I]: A bidirectional codec that can:
//   - Decode input I to type A with validation
//   - Encode type A to output O
//   - Check if a value is of type A
//
// Validation: Decoding returns Either[Errors, A] which represents:
//   - Left(Errors): Validation failed with detailed error information
//   - Right(A): Successfully decoded and validated value
//
// Context: A stack of ContextEntry values that tracks the path through nested structures
// during validation, providing detailed error messages.
//
// # Basic Usage
//
// Creating a simple type:
//
//	nilType := codec.MakeNilType[string]()
//	result := nilType.Decode(nil) // Success
//	result := nilType.Decode("not nil") // Failure
//
// Composing types with Pipe:
//
//	composed := codec.Pipe(typeB)(typeA)
//	// Decodes: I -> A -> B
//	// Encodes: B -> A -> O
//
// # Type Parameters
//
// Most functions use three type parameters:
//   - A: The domain type (the actual Go type being encoded/decoded)
//   - O: The output type for encoding
//   - I: The input type for decoding
//
// # Validation Errors
//
// ValidationError contains:
//   - Value: The actual value that failed validation
//   - Context: The path to the value in nested structures
//   - Message: Human-readable error description
//
// # Integration
//
// This package integrates with:
//   - optics/decoder: For decoding operations
//   - optics/encoder: For encoding operations
//   - either: For validation results
//   - option: For optional type checking
//   - reader: For context-dependent operations
package codec
