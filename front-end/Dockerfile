# Using an official Node.js runtime as a parent image
FROM node:16-alpine

# Copy package.json and package-lock.json (or yarn.lock) to the container
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of your application code to the container
COPY . .

# Build the project in production mode
RUN npm run build && ls -alh

# Install a simple HTTP server to serve the built files
RUN npm install -g serve

# Expose the port that your application will run on
EXPOSE 3000

# Start the server to serve the static files (from dist or build folder)
CMD ["serve", "-s", "dist", "-l", "3000"]
