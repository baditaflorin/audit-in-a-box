export const samplePackageJSON = `{
  "name": "demo-service",
  "version": "1.0.0",
  "dependencies": {
    "express": "4.17.1",
    "lodash": "4.17.20"
  },
  "devDependencies": {
    "jest": "26.6.0"
  }
}`;

export const sampleGoMod = `module example.com/demo

go 1.22

require (
  github.com/go-chi/chi/v5 v5.0.10
  github.com/dgrijalva/jwt-go v3.2.0+incompatible
)`;

export const sampleRequirements = `django==3.2.0
requests==2.25.1
pyyaml==5.3.1`;
