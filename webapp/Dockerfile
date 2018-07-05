FROM node:8

WORKDIR /root
# Copy aha source code
RUN git clone https://github.com/theZiz/aha.git
RUN cd aha; make; make install

# Copy webapp content
RUN install -d /root/webapp
COPY frontend/ /root/webapp/frontend
COPY node_modules/ /root/webapp/node_modules
COPY package-lock.json /root/webapp
COPY scripts/ /root/webapp/scripts
COPY server.js /root/webapp

WORKDIR /root/webapp
CMD node server.js
