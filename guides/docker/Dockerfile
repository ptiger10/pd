FROM yunabe/lgo:latest

# Fetch gopandas
RUN go get -u github.com/ptiger10/pd/... 
RUN lgo installpkg github.com/ptiger10/pd/... 

WORKDIR /notebooks

# To use JupyterLab, replace "notebook" with "lab".
CMD ["jupyter", "notebook", "--ip=0.0.0.0"]
