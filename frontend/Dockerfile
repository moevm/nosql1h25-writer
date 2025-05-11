# Step 1: Build
FROM node:20-alpine3.21 AS builder
WORKDIR /app
COPY ../frontend/package.json package-lock.json ./
RUN npm ci
COPY ../frontend .
RUN npm run build

# Step 2: Run
FROM nginx:1.27.4-alpine3.21
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]