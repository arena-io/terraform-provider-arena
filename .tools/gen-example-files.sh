#!/bin/zsh

tf_files=${PROJECT_ROOT}/internal/provider

for res_file in "${tf_files}"/*_resource.go
do
  if [[ "$res_file" =~ provider\/([^\/]+)_resource.go ]]; then
      domain="arena_${BASH_REMATCH[1]}"
      echo "Domain name: $domain"
      mkdir -p examples/resources/${domain}
      touch examples/resources/${domain}/resource.tf
      cat examples/tg.hcl.tmpl > examples/resources/${domain}/terragrunt.hcl
  else
      echo "Invalid file ${res_file}"
  fi
done

# for ds_file in "${tf_files}"/*_data_source.go
# do
#   if [[ "$ds_file" =~ provider\/([^\/]+)_data_source.go ]]; then
#       domain="arena_${BASH_REMATCH[1]}"
#       echo "Domain name: $domain"
#       mkdir -p examples/data-sources/${domain}
#       touch examples/data-sources/${domain}/data-source.tf
#       cat examples/tg.hcl.tmpl > examples/data-sources/${domain}/terragrunt.hcl
#   else
#       echo "Invalid file ${ds_file}"
#   fi
# done


