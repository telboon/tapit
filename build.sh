#!/bin/bash

cd ./tapit-backend
go build
cd ..

cd ./tapit-frontend
ng build --optimization
cd ..

# copy front-end
cp -r ./tapit-frontend/dist/tapit-frontend/* ./tapit-build/static/
# remove maps
rm ./tapit-build/static/*.map

# copy back-end
cp ./tapit-backend/tapit-backend ./tapit-build/tapit

# run server
./tapit-build/tapit
