require "graphlient"
require "json"

auth = JSON.parse(`./auth/gql-auth`)

client = Graphlient::Client.new('https://gql.twitch.tv/gql',
  headers: {
    "Client-Id": auth["client-id"],
    "Authorization": "OAuth "+auth["token"]
  }
)

result = client.query do
  query do
    currentUser do
      id
      email
    end
  end
end

puts result.data.current_user.id
puts result.data.current_user.email


# update an event!
update_input = {
  id: "YE19n33zQSKNiphX7OYyNw",
  channelID: "5031651",
  title: "Testing 2",
  description: "Testing 2",
  endAt: "2018-08-18T04:00:57.875Z",
  startAt: "2018-08-18T03:00:57.875Z",
  gameID: "488191",
  ownerID: "5031651"
}

result = client.query(input: update_input) do
  mutation(input: :UpdateSingleEventInput!) do
    updateSingleEvent(input: :input) do
      event {
        id
        title
        description
        startAt
        endAt
        game {
          id
          displayName
        }
        channel {
          id
          login
          displayName
        }
        imageURL
      }
    end
  end
end

puts result.data.update_single_event.event.id

# # update an event!
# update_input = {
#   id: "NDfmb3vcQH6Ani0pG9-MhA",
#   channelID: "5031651",
#   title: "Testing 4",
#   description: "Testing 2",
#   endAt: "2018-08-16T04:00:57.875Z",
#   startAt: "2018-08-16T03:00:57.875Z",
#   gameID: "488191",
#   ownerID: "5031651"
# }

# create_input = {
#   channelID: "5031651",
#   title: "Testing 4",
#   description: "Testing 2",
#   endAt: "2018-08-16T04:00:57.875Z",
#   startAt: "2018-08-16T03:00:57.875Z",
#   gameID: "488191",
#   ownerID: "5031651"
# }

# delete_input = {
#   eventID: "cwgIaYqMSwyCeZJs56ylMg",
# }

# # result = client.query(input: delete_input) do
# #   mutation(input: :DeleteEventInput!) do
# #     deleteEvent(input: :input) do
# #       event {
# #         id
# #         title
# #         description
# #         startAt
# #         endAt
# #         game {
# #           id
# #           displayName
# #         }
# #         channel {
# #           id
# #           login
# #           displayName
# #         }
# #         imageURL
# #       }
# #     end
# #   end
# # end


# result = client.query(input: delete_input) do
#     mutation(input: :DeleteEventLeafInput!) do
#       deleteEventLeaf(input: :input) do
#         event {
#           id
#         }
#       end
#     end
#   end

# puts result.data.delete_event_leaf.event.id
# # puts result.data.create_single_event.event.title
# # puts result.data.create_single_event.event.description
# # puts result.data.create_single_event.event.start_at
# # puts result.data.create_single_event.event.end_at