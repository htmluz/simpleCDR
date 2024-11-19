# In development

A rest api to work with freeradius accounting telcobridges ProSBCs CDRs

Steps to get it working

- create a client config with the nastype as other, /etc/freeradius/3.0/clients.conf
- copy the telco dictionary to the dictionary directory and include it on the main dictionary files, normally in /usr/share/freeradius/dictionary
- copy the rest module file to the /etc/freeradius/3.0/mods-available folder
- create a symlink for the rest file in mods-enabled folder
- put the rest module in the default accounting config, /etc/freeradius/3.0/sites-enabled/default
- get running a freeradius
- get running a postgres with the schema provided
