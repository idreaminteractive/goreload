FROM gitpod/workspace-full:2024-05-08-19-39-59

# install air
RUN go install github.com/air-verse/air@latest

# alias all the things
RUN echo 'alias home="cd ${GITPOD_REPO_ROOT}"' | tee -a ~/.bashrc ~/.zshrc