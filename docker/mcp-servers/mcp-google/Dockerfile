###
FROM node:22.12-bullseye AS builder

COPY . /app

WORKDIR /app

RUN --mount=type=cache,target=/root/.npm npm install

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

WORKDIR /app

COPY --from=builder /app/build /app/build
COPY --from=builder /app/package.json /app/package.json
COPY --from=builder /app/package-lock.json /app/package-lock.json

RUN npx playwright install chromium
RUN npm ci --ignore-scripts --omit-dev

CMD ["node", "build/index.js"]