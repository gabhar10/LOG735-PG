FROM node:8

WORKDIR /root
# Copy aha source code
RUN git clone https://github.com/theZiz/aha.git
RUN cd aha; make; make install
# Copy vim source code
RUN git clone https://github.com/vim/vim.git
RUN cd vim; make; make install
RUN ln -s $(which vim) /usr/bin/vi
RUN printf "%s\n%s\n" ":set lines=50 columns=1000" ":syntax on" > /root/.vimrc
# Copy webapp content
RUN install -d /root/webapp
COPY frontend/ /root/webapp/frontend
COPY node_modules/ /root/webapp/node_modules
COPY package-lock.json /root/webapp
COPY scripts/ /root/webapp/scripts
COPY server.js /root/webapp

WORKDIR /root/webapp
CMD node server.js
