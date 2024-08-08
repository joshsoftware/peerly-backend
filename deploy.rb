require 'mina/bundler'
require 'mina/default'
require 'mina/deploy'
require 'mina/git'


set :repository, 'https://github.com/joshsoftware/peerly-backend.git'
set :user, 'ubuntu'
set :forward_agent, true

set :shared_files, [ 
  '.env', 'serviceAccountKey.json'
]

task :staging do 
  set :deploy_to, '/www/peerly-backend'
  set :domain, ENV['PEERLY_STAGING_IP']
  set :branch, 'Deployment'
end

task :production do 
  set :deploy_to, '/www/peerly'
  set :domain, ENV['PEERLY_STAGING_IP']
  set :branch, 'Deployment'
end

task :setup do
   command %{mkdir -p "#{fetch(:deploy_to)}/releases"}
  #  npm install pm2 -g
  #  command %{createdb -U postgres peerly}
  #  command %{pm2 start main --name peerly-backend -- start}
end

task :loadData do
  command %{make seed}
  command %{make loadUser}
end

task :deploy do
  deploy do
    invoke :'git:clone'
    command "git checkout #{fetch(:branch)}"
    command "git pull origin #{fetch(:branch)}"
    invoke :'deploy:link_shared_paths'

    command %{export PATH=$PATH:/usr/local/go/bin}
    command %{go mod tidy}
    command %{go mod vendor}
    command %{go build cmd/main.go}
    # command %{make migrate}
    invoke :'deploy:cleanup' 

    on :launch do
      command %{source ~/.nvm/nvm.sh}
      command %{pm2 restart ~/peerly-ecosystem.config.js}
    end
    
  end
end