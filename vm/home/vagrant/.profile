# ~/.profile: executed by the command interpreter for login shells.
# This file is not read by bash(1), if ~/.bash_profile or ~/.bash_login
# exists.
# see /usr/share/doc/bash/examples/startup-files for examples.
# the files are located in the bash-doc package.

# the default umask is set in /etc/profile; for setting the umask
# for ssh logins, install and configure the libpam-umask package.
#umask 022

# if running bash
if [ -n "$BASH_VERSION" ]; then
    # include .bashrc if it exists
    if [ -f "$HOME/.bashrc" ]; then
        . "$HOME/.bashrc"
    fi
fi

# set PATH so it includes user's private bin if it exists
if [ -d "$HOME/bin" ] ; then
    PATH="$HOME/bin:$PATH"
fi

export PS1="\\[\e[33;1m\]\u\[\e[36;1m\]@\[\e[32;1m\]\H \[\e[36;1m\]|\d \t| \e[33;1m\]\w \e[37;1m\]\n\$> \[\e[0m\]"

#############################
# cd
#############################
alias ..='cd ..'
alias cdproj='cd /vagrant'

#############################
# df & du
#############################
alias dfh='df -h'
alias duh='du -h'

#############################
# Profile
#############################
alias viprof='vim ~/.profile'
alias scprof='source ~/.profile'

#############
# ls
#############
function ll(){
	ls -lah $1
}

#############
# Nginx
#############
alias nginxRel='sudo systemctl reload nginx'

###########
# shutdown
###########
alias shutup='sudo shutdown -P -h now'
