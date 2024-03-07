# topb

**topb** is a tool for automatically generating conversion methods from Go structs to Protobuf message types. It parses Go source files and generates corresponding `ToPb` methods, simplifying the process of manually writing these conversion methods.

## Installation

\`\`\`bash
go install github.com/your-username/topb@latest
\`\`\`

## Usage

In your Go source file, add the following comments for the structs that need to generate `ToPb` methods:

\`\`\`go
// gen:topb
//
//go:generate topb -in user.go
type User struct {
...
}
\`\`\`

Then, run the following command in the terminal:

\`\`\`bash
go generate ./...
\`\`\`

**topb** will automatically generate a file named `autogen_topb_user.go`, which contains the `ToPb` method for the `User` struct.

## Example

Suppose you have the following struct definition:

\`\`\`go
// gen:topb
//
//go:generate topb -in user.go
type User struct {
Model

    ID        uint32 \`protobuf:"varint,1,opt,name=id" json:"id" gorm:"index;comment:'User ID'"\`
    Name      string \`protobuf:"bytes,2,opt,name=name" json:"name" gorm:"comment:'User name'"\`
    Email     string \`protobuf:"-" json:"email" gorm:"comment:'User email'"\`
}
\`\`\`

**topb** will generate the following content:

\`\`\`go
// Code generated by topb; DO NOT EDIT.

package main

import "github.com/tikbox/topb/pb"

func (m *User) ToPb() *pb.User {
return &pb.User{
Id:   m.ID,
Name: m.Name,
}
}
\`\`\`

As you can see, the generated `ToPb` method converts the `User` struct to the corresponding Protobuf message type `pb.User`.

## Notes

- **topb** only generates `ToPb` methods and does not generate reverse conversion methods from Protobuf message types to Go structs.
- For embedded fields in the struct, **topb** will skip them and not generate conversion code for them.
- If a struct field uses the `protobuf:"-"` tag, **topb** will ignore that field.

## Contribution

Feel free to submit issues or Pull Requests to improve **topb**. If you have any questions or suggestions, you can also contact us at any time.

## Advanced Usage

In addition to adding comments in the source file to specify the structs for which to generate `ToPb` methods, **topb** also supports some other ways.

### Command Line Options

You can use command line options to specify the files or directories to process:

\`\`\`bash
topb -file=path/to/file.go
topb -dir=path/to/dir
\`\`\`

Using the `-file` option specifies a single file, and using the `-dir` option specifies an entire directory. **topb** will parse all Go files in these files or directories and generate `ToPb` methods for eligible structs.

### Exclude Files

If you don't want to generate `ToPb` methods for certain files, you can use the `exclude` tag to exclude them:

\`\`\`go
// gen:topb
// gen:topb,exclude
type User struct {
...
}
\`\`\`

In this example, **topb** will skip the `User` struct and not generate a `ToPb` method for it.

### Custom Import Path

By default, **topb** will import the `"github.com/tikbox/topb/pb"` package. If you are using another Protobuf package, you can specify the import path by adding a comment:

\`\`\`go
// gen:topb
// gen:topb,import="your/custom/pbpkg"
type User struct {
...
}
\`\`\`

In this example, **topb** will import the `"your/custom/pbpkg"` package and use it to generate the `ToPb` method.

## Testing

**topb** includes some test cases for verifying the correctness of the generated code. You can run the tests with the following command:

\`\`\`bash
go test ./...
\`\`\`

If you modify the **topb** code, please make sure to add or modify the corresponding test cases to ensure code correctness.

## Roadmap

Here are some planned features and improvements for the **topb** project:

- Support generating reverse conversion methods from Protobuf message types to Go structs
- Support customizing the generated method names
- Support filtering out specific struct fields
- Improve the performance and scalability of the code generator
- Add more test cases to increase code coverage

If you have any other suggestions or requirements, feel free to bring them up for discussion.

## Contributors

Thanks to the following people for contributing to the **topb** project:

- [Your Name](https://github.com/your-username)

## License

**topb** is licensed under the [MIT License](LICENSE).

Welcome to use, contribute, and share **topb**!
