###
FROM node:22.12-bullseye AS builder

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package.json package-lock.json ./

# Install dependencies without running scripts to prevent automatic build
RUN npm install --ignore-scripts

# Copy the entire source directory
COPY src ./src
COPY tsconfig.json ./

# Build the project
RUN npm run build

###
FROM node:22.12-bullseye AS release

ENV NODE_ENV=production

# install chromium and dependencies for playwright
RUN apt-get update && apt-get install -y --no-install-recommends \                                
	chromium \                                                                                        
	libasound2 \                                                                                      
	libgbm-dev \                                                                                      
	libgtk-3-0 \                                                                          
	libnss3 \                                                                                         
	libx11-xcb1 \                                                                                     
	libxcomposite1 \                                                                                  
	libxdamage1 \                                                                                     
	libxext6 \                                                                                        
	libxfixes3 \                                                                                      
	libxi6 \                                                                                          
	libxkbcommon0 \                                                                                   
	libxrandr2 \                                                                                      
	libxrender1 \                                                                                     
	fonts-dejavu \                                                                                    
	&& rm -rf /var/lib/apt/lists/*                                                                    


# Set the working directory
WORKDIR /app

# Copy the built files from the builder stage
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/package-lock.json ./package-lock.json

RUN npx playwright install chromium

RUN npm ci --ignore-scripts --omit=dev

CMD ["node", "dist/index.js"]