require 'mina/bundler'
require 'mina/default'
require 'mina/deploy'
require 'mina/git'

set :repository, 'https://github.com/joshsoftware/peerly-backend.git'
set :user, 'ubuntu'
set :shared_files, [ 
  '.env', 'serviceAccountKey.json'
]

# Define environments
task :staging do 
  set :deploy_to, '/www/peerly-backend'
  set :domain, 'pg-stage-intranet.joshsoftware.com'
  set :branch, 'Deployment'
end

task :production do 
  set :deploy_to, '/www/peerly-backend'
  set :domain, 'intranet.joshsoftware.com'
  set :branch, 'main'
end

task :setup do
   command %{mkdir -p "#{fetch(:deploy_to)}/releases"}
end

task :loadData do
  command %{make seed}
  command %{make loadUser}
end

# Main deploy task
task :deploy do
  deploy do
    invoke :'git:clone'
    invoke :'deploy:link_shared_paths'
    command %{export PATH=$PATH:/usr/local/go/bin}
    command %{go mod tidy}
    command %{go mod vendor}
    command %{go build cmd/main.go}
    on :launch do
      command %{source ~/.nvm/nvm.sh}
      command %{pm2 restart ~/peerly-ecosystem.config.js}
    end
  end
end