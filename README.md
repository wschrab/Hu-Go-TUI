# Hu Go TUI


**NOTE:**

I have only tried this with the papermod theme
This codebase is still a work in progress and will likely not work for you without some editing. 
Most importantly in the model.go file there is a line 'testDir := `C:\Users\wdsch\second-portfolio`'
This must be edited to be the base of where your hugo file is stored on your computer.

## Required HUGO file structure
Your file structure, for posts at least, should be in the following format:

- content   
    - posts
        - _index.yaml
        - post1
            - index.yaml
            - media1.jpg
            - etc media
        - post2
            - index.yaml
            - media1.jpg
            - etc media
        - etc posts