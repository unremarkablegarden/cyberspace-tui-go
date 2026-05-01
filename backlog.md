# Backlog

Stuff found while working on the code

## Stuff done
- [X] Post Details does not know where to return, always return to main feed
- [X] Profile does not know where to return, always return to main feed
- [X] No username when post model added
- [X] Topics shows 0 posts counter for all topics
- [X] When open notifications panel, post id is always empty
- [X] Same error as empty post id in bookmarks
- [X] When theme changes, no file is being save to persist changes
- [X] Logout functionality not implemented
- [X] Some list paginated events uses similar logic, maybe a generic function can help with that

## Stuff to fix
- [ ] With current changes, each screen change triggers a new request (local cache?)
- [ ] When a post is saved in the post detail view, there is no "unsave" functionality
- [ ] Update help keys
- [ ] Note composer "IsEdit" flag works only to edit a note, if "IsEdit" is false, falls back to "new note" instead of open a note in read mode or something
- [ ] Only first pagination is working

- [ ] Move *_items to be more clear

## Stuff to do
- [ ] Menu for easy navigation
- [ ] Support images in terminals that support _TGP_
- [ ] Expand button in feed screen
