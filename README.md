# cloud-paint

This project only exists because I wanted to learn the go programming language. And because I am also a big fan of big picture diagrams for software I created [cloudpaint](https://github.com/nrekretep/cloudpaint) and [cloudpaint-cli](https://github.com/nrekretep/cloudpaint-cli).

## What is this tool good for?

Cloud-paint is a tool for cloudfoundry to visualize an app running on cloud-foundry with all its dependencies to buildpacks, stacks, services, network policies and much more. 

The target audience includes developers and architects which want to check if the running app really conforms to their architecture guidelines. 
Also platform teams can benefit from this tool. 

## What is not included? 

cloud-paint currently only supports plantuml diagrams. From cloud-paint you can only get the plain text form of the diagram in plantuml syntax. You need to send this raw diagram to a plantuml renderer of your choice. 

# Documentation

## Project documentation

Additional documentation for the project as a whole can be found in the [docs](docs/) directory.

## Package documentation

Every package includes a doc.go file containing only a package statement and some documentation for this package.

See also [Godoc: documenting Go code](https://blog.golang.org/godoc-documenting-go-code).

# Testing

Well I created this project to get going with go. And the first time I just spent some time learning go basics. But as a professional software engineer I know that test driven development is also included in the first steps of learning any new programming language. 

So I will try to stick to a strict red-green-refactor cycle. 

To create the tests I will use [goconvey](http://goconvey.co). Goconvey is easy to setup and use. [Just watch this video and you are good to go](https://www.youtube.com/watch?v=wlUKRxWEELU). 

