
# WebCSV - displaying csv files in your browser

This web application provides a user-friendly interface for browsing and visualizing CSV files within a specified directory. Upon requests, the application scans the designated directory, identifies all files with the '.csv' extension, and presents them as a list in the web interface. Users can then select a CSV file from the list, and the application will dynamically load and display the file's contents in a tabular format within the same web page using JavaScript.

 It mainly use  [git@github.com>:derekeder/csv-to-html-table.git](git@github.com>:derekeder/csv-to-html-table.git)

common usage, for me is :

- data is shared from mounted volume in docker
- volume is syncthinc'ed data from windows,
- https service is proxied via prx.si to dev host.

Base URL is :
  <http://fqdn.host.dom:8080/list>
