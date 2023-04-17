FROM debian:stable-slim

COPY ./dist/ /app/
COPY ./server/.assets/emacs/htmlize.el /usr/share/emacs/site-lisp/
COPY ./server/.assets/emacs/ox-gfm.el /usr/share/emacs/site-lisp/

WORKDIR "/app"

RUN apt-get update > /dev/null && \
    apt-get install -y --no-install-recommends apt-utils && \
    apt-get install -y curl emacs-nox ffmpeg zip poppler-utils > /dev/null && \
    # org-mode: pdf export (with a light latex distribution)
    cd && apt-get install -y wget perl > /dev/null && \
    export CTAN_REPO="http://mirror.las.iastate.edu/tex-archive/systems/texlive/tlnet" && \
    curl -sL "https://yihui.name/gh/tinytex/tools/install-unx.sh" | sh && \
    mv ~/.TinyTeX /usr/share/tinytex && \
    /usr/share/tinytex/bin/x86_64-linux/tlmgr install wasy && \
    /usr/share/tinytex/bin/x86_64-linux/tlmgr install ulem && \
    /usr/share/tinytex/bin/x86_64-linux/tlmgr install marvosym && \
    /usr/share/tinytex/bin/x86_64-linux/tlmgr install wasysym && \
    /usr/share/tinytex/bin/x86_64-linux/tlmgr install xcolor && \
    /usr/share/tinytex/bin/x86_64-linux/tlmgr install listings && \
    /usr/share/tinytex/bin/x86_64-linux/tlmgr install parskip && \
    /usr/share/tinytex/bin/x86_64-linux/tlmgr install float && \
    /usr/share/tinytex/bin/x86_64-linux/tlmgr install wrapfig && \
    /usr/share/tinytex/bin/x86_64-linux/tlmgr install sectsty && \
    /usr/share/tinytex/bin/x86_64-linux/tlmgr install capt-of && \
    /usr/share/tinytex/bin/x86_64-linux/tlmgr install epstopdf-pkg && \
    /usr/share/tinytex/bin/x86_64-linux/tlmgr install cm-super && \
    ln -s /usr/share/tinytex/bin/x86_64-linux/pdflatex /usr/local/bin/pdflatex && \
    apt-get purge -y --auto-remove perl wget && \
    # Cleanup
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /tmp/* && \
    find /usr/share/ -name 'doc' | xargs rm -rf && \
    find /usr/share/emacs -name '*.pbm' | xargs rm -f && \
    find /usr/share/emacs -name '*.png' | xargs rm -f && \
    find /usr/share/emacs -name '*.xpm' | xargs rm -f

RUN useradd filestash && \
    chown -R filestash:filestash /app/

USER filestash

EXPOSE 8334

CMD ["/app/filestash"]
