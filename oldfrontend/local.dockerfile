FROM node:20.11.1-alpine

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /app

# Copy package.json and pnpm-lock.yaml
COPY package.json pnpm-lock.yaml* ./

# Install pnpm and dependencies
RUN npm install -g pnpm && pnpm install

# Copy the rest of the application code
COPY . .

# Set permissions
RUN chown -R appuser:appgroup /app
RUN chmod -R u+w /app


# Switch to the non-root user
USER appuser

# Expose the port
EXPOSE 5173

# Start the development server
CMD ["pnpm", "run", "dev", "--host"]
